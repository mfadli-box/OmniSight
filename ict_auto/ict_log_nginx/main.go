package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	PostgresConn       string
	ElasticsearchURL   string
	HttpFloodThreshold int
}

type NginxLog struct {
	Timestamp     string      `json:"@timestamp"`
	Host          string      `json:"host"`
	ServerIP      string      `json:"server_ip"`
	ClientIP      string      `json:"client_ip"`
	Xff           string      `json:"xff"`
	Domain        string      `json:"domain"`
	URL           string      `json:"url"`
	Referer       string      `json:"referer"`
	Args          string      `json:"args"`
	UpstreamTime  json.Number `json:"upstreamtime"`
	ResponseTime  json.Number `json:"responsetime"`
	RequestMethod string      `json:"request_method"`
	Status        json.Number `json:"status"`
	Size          json.Number `json:"size"`
	RequestBody   string      `json:"request_body"`
	RequestLength string      `json:"request_length"`
	Protocol      string      `json:"protocol"`
	UpstreamHost  string      `json:"upstreamhost"`
	FileDir       string      `json:"file_dir"`
	UserAgent     string      `json:"http_user_agent"`
	GeoIP         *struct {
		Geo *struct {
			CountryISOCode string `json:"country_iso_code"`
		} `json:"geo"`
	} `json:"geoip"`
}

type EsHit struct {
	ID     string   `json:"_id"`
	Source NginxLog `json:"_source"`
}

type EsSearchResponse struct {
	Hits struct {
		Hits []EsHit `json:"hits"`
	} `json:"hits"`
}

type AttackKey struct {
	ClientIP    string
	TrafficType string
	Domain      string
}

type BypassRule struct {
	Domain      string
	URLPath     string
	ArgsPattern string
}

type Cache struct {
	mu          sync.RWMutex
	exactIPs    map[string]bool
	cidrNets    []*net.IPNet
	banned      map[string]bool
	bypassRules []BypassRule
	loadedAt    time.Time
}

var (
	CacheInstance = &Cache{}
	PgSQL         *sql.DB
)

const cacheRefreshInterval = 5 * time.Minute

func (c *Cache) Refresh(db *sql.DB) {
	c.mu.Lock()
	defer c.mu.Unlock()

	exactIPs := make(map[string]bool)
	var cidrNets []*net.IPNet

	rows, err := db.Query(`
		SELECT ip_or_cidr FROM "ict_ip_whitelist"
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var entry string
			if err := rows.Scan(&entry); err != nil {
				continue
			}
			entry = strings.TrimSpace(entry)
			if strings.Contains(entry, "/") {
				_, ipNet, err := net.ParseCIDR(entry)
				if err == nil {
					cidrNets = append(cidrNets, ipNet)
				}
			} else if ip := net.ParseIP(entry); ip != nil {
				exactIPs[entry] = true
			}
		}
	}
	if rows != nil {
		rows.Close()
	}

	banned := make(map[string]bool)
	rowsBan, err := db.Query(`
		SELECT ip FROM "ict_ip_blacklist"
		WHERE  expires_at IS NULL OR expires_at > NOW()
	`)
	if err == nil {
		defer rowsBan.Close()
		for rowsBan.Next() {
			var ip string
			if err := rowsBan.Scan(&ip); err == nil {
				banned[ip] = true
			}
		}
	}
	if rowsBan != nil {
		rowsBan.Close()
	}

	var bypassRules []BypassRule
	rowsBypass, err := db.Query(`
		SELECT domain, url_path, COALESCE(args_pattern, '')
		FROM   "ict_waf_bypass_rule"
	`)
	if err == nil {
		defer rowsBypass.Close()
		for rowsBypass.Next() {
			var r BypassRule
			if err := rowsBypass.Scan(&r.Domain, &r.URLPath, &r.ArgsPattern); err == nil {
				bypassRules = append(bypassRules, r)
			}
		}
	}
	if rowsBypass != nil {
		rowsBypass.Close()
	}

	c.exactIPs = exactIPs
	c.cidrNets = cidrNets
	c.banned = banned
	c.bypassRules = bypassRules
	c.loadedAt = time.Now()
	log.Printf("Cache refreshed: %d whitelist IPs, %d CIDR, %d banned, %d bypass rules",
		len(exactIPs), len(cidrNets), len(banned), len(bypassRules))
}

func (c *Cache) IsWhitelisted(ipStr string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	if c.exactIPs[ipStr] {
		return true
	}
	for _, cidr := range c.cidrNets {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func (c *Cache) IsBanned(ipStr string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.banned[ipStr]
}

func (c *Cache) AddBanned(ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.banned[ip] = true
}

func (c *Cache) IsBypassRule(domain, urlPath, args string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, r := range c.bypassRules {
		if (r.Domain == "*" || r.Domain == domain) && r.URLPath == urlPath {
			if r.ArgsPattern == "" || args == "" || strings.Contains(args, r.ArgsPattern) {
				return true
			}
		}
	}
	return false
}

func drainBody(body io.ReadCloser) {
	if body == nil {
		return
	}
	defer body.Close()
	io.Copy(io.Discard, body)
}

type IPTracker struct {
	mu    sync.RWMutex
	count map[string]int
}

func (t *IPTracker) Incr(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count[key]++
}

func (t *IPTracker) Get(key string) int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.count[key]
}

func (t *IPTracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count = make(map[string]int)
}

var TrackerInstance = &IPTracker{count: make(map[string]int)}

func LoadConfig() *Config {
	godotenv.Load()
	thresh := 100
	if envThresh := os.Getenv("FT_HTTP"); envThresh != "" {
		if val, err := strconv.Atoi(envThresh); err == nil {
			thresh = val
		}
	}
	PG_Host := os.Getenv("PG_HOST")
	PG_Port := os.Getenv("PG_PORT")
	PG_User := os.Getenv("PG_USER")
	PG_Pass := os.Getenv("PG_PASS")
	PG_Data := os.Getenv("PG_DATA")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		PG_User, PG_Pass, PG_Host, PG_Port, PG_Data)
	return &Config{
		PostgresConn:       dsn,
		ElasticsearchURL:   os.Getenv("ES_LINK"),
		HttpFloodThreshold: thresh,
	}
}

func InitDB(connStr string) {
	var err error
	PgSQL, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	PgSQL.SetMaxOpenConns(25)
	PgSQL.SetMaxIdleConns(10)
	PgSQL.SetConnMaxLifetime(5 * time.Minute)
	if err = PgSQL.Ping(); err != nil {
		log.Fatalf("Database did not respond to Ping: %v", err)
	}
}

func initElasticClient() (*elasticsearch.Client, error) {
	elasticAddr := os.Getenv("ES_LINK")
	if elasticAddr == "" {
		elasticAddr = "http://elasticsearch:9200"
	}
	customTransport := &http.Transport{
		DisableKeepAlives:   false,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	cfg := elasticsearch.Config{
		Addresses: []string{elasticAddr},
		Transport: customTransport,
	}
	return elasticsearch.NewClient(cfg)
}

var (
	RegexSQLI           = regexp.MustCompile(`(?i)(UNION|SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|WHERE|OR\s+\d=\d|--|#|\/\*|AND\s+\d=\d|ORDER\s+BY|GROUP\s+BY)`)
	RegexXSS            = regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onload=|alert\(|%3Cscript|<svg|<iframe|confirm\(|prompt\()`)
	RegexTraversal      = regexp.MustCompile(`(?i)(\.\.\/|\.\.\\|%2e%2e%2f|etc\/passwd|boot\.ini|win\.ini|%5c\.\.)`)
	RegexRCE            = regexp.MustCompile(`(?i)(bin\/sh|bin\/bash|cmd\.exe|powershell|wget\s|curl\s|eval\(|passthru|shell_exec|system\(|popen\()`)
	RegexSensitiveFiles = regexp.MustCompile(`(?i)(\.env|\.git|\.docker|config\.php|wp-config\.php|db\.php|\.bak|\.sql|\.yaml)`)
	RegexBotScanner     = regexp.MustCompile(`(?i)(nikto|sqlmap|dirbuster|w3af|nmap|acunetix|masscan|python-requests|curl|hydra|gobuster|wfuzz|amass|zgrab)`)
)

func ClassifyTraffic(url, args, body, ua string) string {
	payload := strings.ToLower(url + " " + args + " " + body)
	uaLower := strings.ToLower(ua)
	if RegexBotScanner.MatchString(uaLower) {
		return "BOT_SCANNER"
	}
	if RegexSQLI.MatchString(payload) {
		return "SQL_INJECTION"
	}
	if RegexXSS.MatchString(payload) {
		return "XSS"
	}
	if RegexTraversal.MatchString(payload) {
		return "PATH_TRAVERSAL"
	}
	if RegexRCE.MatchString(payload) {
		return "RCE_COMMAND_INJECTION"
	}
	if RegexSensitiveFiles.MatchString(payload) {
		return "SENSITIVE_FILE_PROBING"
	}
	return "NORMAL"
}

func GetThreatWeight(trafficType string) int {
	switch trafficType {
	case "RCE_COMMAND_INJECTION":
		return 50
	case "SQL_INJECTION":
		return 40
	case "PATH_TRAVERSAL":
		return 35
	case "XSS":
		return 20
	case "SENSITIVE_FILE_PROBING":
		return 15
	case "BOT_SCANNER":
		return 10
	default:
		return 0
	}
}

func UpdateThreatScoreDB(db *sql.DB, ip string, trafficType string) (int, bool) {
	if CacheInstance.IsWhitelisted(ip) {
		return 0, false
	}
	if CacheInstance.IsBanned(ip) {
		return 150, true
	}
	weight := GetThreatWeight(trafficType)
	if weight == 0 {
		return TrackerInstance.Get(ip), TrackerInstance.Get(ip) >= 100
	}
	TrackerInstance.mu.Lock()
	TrackerInstance.count[ip] += weight
	finalScore := TrackerInstance.count[ip]
	if finalScore > 150 {
		TrackerInstance.count[ip] = 150
		finalScore = 150
	}
	TrackerInstance.mu.Unlock()
	if finalScore >= 100 {
		expiryTime := time.Now().Add(24 * time.Hour)
		query := `
			INSERT INTO "ict_ip_blacklist" (
				id, ip, threat_score, reason, expires_at
			) VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (ip) DO UPDATE
			 SET threat_score = EXCLUDED.threat_score, expires_at = EXCLUDED.expires_at;
		`
		newUUID := uuid.NewString()
		_, err := db.Exec(query, newUUID, ip, finalScore, trafficType, expiryTime)
		if err != nil {
			log.Printf("Failed to save IP block %s to database: %v", ip, err)
		}
		CacheInstance.AddBanned(ip)
		return finalScore, true
	}
	return finalScore, false
}

func StartSyncWorker(es *elasticsearch.Client, threshold int) {
	runSyncTask(es, threshold)
	ticker := time.NewTicker(1 * time.Minute)
	cacheTicker := time.NewTicker(cacheRefreshInterval)
	for {
		select {
		case <-ticker.C:
			runSyncTask(es, threshold)
		case <-cacheTicker.C:
			CacheInstance.Refresh(PgSQL)
		}
	}
}

func runSyncTask(es *elasticsearch.Client, threshold int) {
	ctx := context.Background()
	now := time.Now()
	eod := now.AddDate(0, 0, -1)
	beforeIndex := fmt.Sprintf("logstash_%s", eod.Format("2006.01.02"))
	currentIndex := fmt.Sprintf("logstash_%s", now.Format("2006.01.02"))
	syncAndAnalyze(ctx, es, beforeIndex, eod, threshold)
	syncAndAnalyze(ctx, es, currentIndex, now, threshold)
	cleanupOldIndices(ctx, es, eod)
}

func restartSelf() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to find binary path: %v", err)
	}
	cmd := exec.Command(self, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to restart application: %v", err)
	}
	os.Exit(0)
}

func syncAndAnalyze(ctx context.Context, es *elasticsearch.Client, indexName string, now time.Time, threshold int) {
	cleanJsonQuery := `{"query":{"bool":{"must":{"match":{"error.message":"Error decoding JSON: invalid character 'x' in string escape code"}}}}}`
	res, err := es.DeleteByQuery([]string{indexName}, strings.NewReader(cleanJsonQuery),
		es.DeleteByQuery.WithContext(ctx), es.DeleteByQuery.WithConflicts("proceed"))
	drainBody(res.Body)

	cleanStatus0 := `{"query":{"match":{"status":"0"}}}`
	res, err = es.DeleteByQuery([]string{indexName}, strings.NewReader(cleanStatus0),
		es.DeleteByQuery.WithContext(ctx), es.DeleteByQuery.WithConflicts("proceed"))
	drainBody(res.Body)

	cleanIpQuery := `{"query":{"bool":{"should":[{"term":{"client_ip":""}},{"bool":{"must_not":{"exists":{"field":"client_ip"}}}}],"minimum_should_match":1}}}`
	res, err = es.DeleteByQuery([]string{indexName}, strings.NewReader(cleanIpQuery),
		es.DeleteByQuery.WithContext(ctx), es.DeleteByQuery.WithConflicts("proceed"))
	drainBody(res.Body)

	queryFindDates := `
		SELECT DISTINCT date_str FROM (
			SELECT DISTINCT TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_str
			FROM   "ict_nginx_log"
			WHERE  client_ip = ''
			UNION
			SELECT DISTINCT TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_str
			FROM   "ict_nginx_app"
			WHERE  client_ip = ''
			UNION
			SELECT DISTINCT TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_str
			FROM   "ict_nginx_atc"
			WHERE  client_ip = ''
		) t WHERE date_str IS NOT NULL
	`
	rowsDates, errDates := PgSQL.QueryContext(ctx, queryFindDates)
	if errDates == nil {
		var affectedDates []string
		for rowsDates.Next() {
			var d string
			if err := rowsDates.Scan(&d); err == nil {
				affectedDates = append(affectedDates, d)
			}
		}
		rowsDates.Close()
		if len(affectedDates) > 0 {
			log.Printf("Log %s Cleaned client_ip - %d", indexName, len(affectedDates))
			txClean, errTxClean := PgSQL.BeginTx(ctx, nil)
			if errTxClean == nil {
				cleanSuccess := true
				for _, dateStr := range affectedDates {
					_, err1 := txClean.ExecContext(ctx, `
						DELETE FROM "ict_nginx_log"
						WHERE  client_ip = '' AND TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
					`, dateStr)
					_, err2 := txClean.ExecContext(ctx, `
						DELETE FROM "ict_nginx_app"
						WHERE  client_ip = '' AND TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
					`, dateStr)
					_, err3 := txClean.ExecContext(ctx, `
						DELETE FROM "ict_nginx_atc"
						WHERE  client_ip = '' AND TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
					`, dateStr)
					_, err4 := txClean.ExecContext(ctx, `
						DELETE FROM "ict_nginx_atc_sum"
						WHERE  client_ip = '' AND date = $1
					`, dateStr)
					if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
						continue
					}
					queryUpdateSLA := `
						WITH metrics AS (
							SELECT
								COUNT(*) AS total,
								COUNT(CASE WHEN (status::int >= 200 AND status::int < 300) AND (traffic_type = 'NORMAL' OR traffic_type = 'WHITELISTED_TRAFFIC') THEN 1 END) AS success,
								COUNT(CASE WHEN status::int >= 400 AND status::int < 500 THEN 1 END) AS client_err,
								COUNT(CASE WHEN status::int >= 500 THEN 1 END) AS server_err,
								COALESCE(AVG(NULLIF(responsetime, '')::numeric), 0.0000) AS avg_time
							FROM (
								SELECT timestamp, status, traffic_type, responsetime
								FROM   "ict_nginx_log"
								WHERE  TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
								UNION ALL
								SELECT timestamp, status, traffic_type, responsetime
								FROM   "ict_nginx_app"
								WHERE  TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
								UNION ALL
								SELECT timestamp, status, traffic_type, responsetime
								FROM   "ict_nginx_atc"
								WHERE  TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
							) combined
						),
						attack_metrics AS (
							SELECT COUNT(*) AS total_attacks
							FROM   "ict_nginx_atc"
							WHERE  TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
						)
						UPDATE "ict_nginx_sla" s
						SET total_requests = m.total,
							successful_requests = m.success,
							client_errors = m.client_err,
							server_errors = m.server_err,
							attack_requests = a.total_attacks,
							avg_response_time = m.avg_time,
							sla_percentage = CASE WHEN m.total > 0 THEN (m.success::numeric / m.total::numeric) * 100 ELSE 0.00 END
						FROM metrics m, attack_metrics a WHERE s.date = $1
					`
					if _, errSLA := txClean.ExecContext(ctx, queryUpdateSLA, dateStr); errSLA != nil {
						cleanSuccess = false
						break
					}
				}
				if cleanSuccess {
					_ = txClean.Commit()
					log.Printf("Log %s SLA updated.", indexName)
				} else {
					txClean.Rollback()
				}
			}
		}
	}
	if rowsDates != nil {
		rowsDates.Close()
	}

	query := `{"size": 2500, "query": {"match_all": {}}}`
	resSearch, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(strings.NewReader(query)))
	if err != nil {
		restartSelf()
		return
	}
	defer drainBody(resSearch.Body)
	if resSearch.IsError() {
		log.Printf("Log %s Kosong", indexName)
		return
	}
	var searchResult EsSearchResponse
	if err := json.NewDecoder(resSearch.Body).Decode(&searchResult); err != nil {
		log.Printf("Failed to decode search result: %v", err)
		return
	}
	if len(searchResult.Hits.Hits) == 0 {
		return
	}
	TrackerInstance.Reset()
	for _, hit := range searchResult.Hits.Hits {
		TrackerInstance.Incr(hit.Source.ClientIP)
	}

	var batchTotal, batchSuccess, batchClientErr, batchServerErr, batchAttack int64
	var totalResponseTime float64
	batchAttackSummary := make(map[AttackKey]int64)
	var deletedDocIDs []string

	tx, err := PgSQL.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	txOpened := true
	defer func() {
		if txOpened {
			tx.Rollback()
		}
	}()

	stmtNormal, err := tx.PrepareContext(ctx, `
		INSERT INTO "ict_nginx_log" (
			id, timestamp, host, server_ip, client_ip, country_iso, xff,
			domain, url, referer, args, upstreamtime, responsetime,
			request_method, status, size, request_body, request_length,
			protocol, upstreamhost, file_dir, http_user_agent, traffic_type
		) VALUES (
		 	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		)`)
	if err != nil {
		log.Printf("Failed to prepare normal insert: %v", err)
		return
	}
	defer stmtNormal.Close()

	stmtAttack, err := tx.PrepareContext(ctx, `
		INSERT INTO "ict_nginx_atc" (
			id, timestamp, host, server_ip, client_ip, country_iso, xff,
			domain, url, referer, args, upstreamtime, responsetime,
			request_method, status, size, request_body, request_length,
			protocol, upstreamhost, file_dir, http_user_agent, traffic_type
		) VALUES (
		 	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		)`)
	if err != nil {
		log.Printf("Failed to prepare attack insert: %v", err)
		return
	}
	defer stmtAttack.Close()

	stmtApp, err := tx.PrepareContext(ctx, `
		INSERT INTO "ict_nginx_app" (
			id, timestamp, host, server_ip, client_ip, country_iso, xff,
			domain, url, referer, args, upstreamtime, responsetime,
			request_method, status, size, request_body, request_length,
			protocol, upstreamhost, file_dir, http_user_agent, traffic_type
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23
		)`)
	if err != nil {
		log.Printf("Failed to prepare app insert: %v", err)
		return
	}
	defer stmtApp.Close()

	for _, hit := range searchResult.Hits.Hits {
		logData := hit.Source

		if strings.TrimSpace(logData.ClientIP) == "" || logData.ClientIP == "-" {
			continue
		}
		if strings.TrimSpace(logData.Domain) == "" || logData.Domain == "-" {
			continue
		}

		statusInt64, _ := logData.Status.Int64()
		statusInt := int(statusInt64)
		statusStr := logData.Status.String()
		responseTimeStr := logData.ResponseTime.String()
		upstreamTimeStr := logData.UpstreamTime.String()
		sizeStr := logData.Size.String()
		countryIso := "-"
		if logData.GeoIP != nil && logData.GeoIP.Geo != nil && logData.GeoIP.Geo.CountryISOCode != "" {
			countryIso = logData.GeoIP.Geo.CountryISOCode
		}
		reqLenInt, err_int := strconv.ParseInt(logData.RequestLength, 10, 64)
		var requestLengthParam interface{}
		if err_int != nil || reqLenInt <= 0 {
			requestLengthParam = sql.NullInt64{Valid: false}
		} else {
			requestLengthParam = reqLenInt
		}
		respTimeFloat, _ := logData.ResponseTime.Float64()
		totalResponseTime += respTimeFloat

		totalHitByIP := TrackerInstance.Get(logData.ClientIP)
		var trafficType string
		if totalHitByIP > threshold {
			trafficType = "HTTP_FLOOD"
		} else {
			trafficType = ClassifyTraffic(
				logData.URL, logData.Args, logData.RequestBody, logData.UserAgent,
			)
		}
		batchTotal++

		isClientWhitelisted := CacheInstance.IsWhitelisted(logData.ClientIP)
		isRuleBypassed := CacheInstance.IsBypassRule(logData.Domain, logData.URL, logData.Args)
		isAttack := trafficType != "NORMAL"

		var targetStmt *sql.Stmt

		if isAttack {
			if isClientWhitelisted {
				targetStmt = stmtApp
				trafficType = "WH_" + trafficType
				if statusInt >= 200 && statusInt < 300 {
					batchSuccess++
				} else if statusInt >= 400 && statusInt < 500 {
					batchClientErr++
				} else if statusInt >= 500 {
					batchServerErr++
				}
			} else if isRuleBypassed {
				targetStmt = stmtApp
				trafficType = "BR_" + trafficType
				if statusInt >= 200 && statusInt < 300 {
					batchSuccess++
				} else if statusInt >= 400 && statusInt < 500 {
					batchClientErr++
				} else if statusInt >= 500 {
					batchServerErr++
				}
			} else {
				targetStmt = stmtAttack
				batchAttack++
				if statusInt >= 400 && statusInt < 500 {
					batchClientErr++
				} else if statusInt >= 500 {
					batchServerErr++
				}
				key := AttackKey{
					ClientIP:    logData.ClientIP,
					TrafficType: trafficType,
					Domain:      logData.Domain,
				}
				batchAttackSummary[key]++
				UpdateThreatScoreDB(PgSQL, logData.ClientIP, trafficType)
			}
		} else {
			if isClientWhitelisted {
				targetStmt = stmtApp
				trafficType = "WH_" + trafficType
			} else if isRuleBypassed {
				targetStmt = stmtNormal
				trafficType = "NORMAL_RULE"
			} else {
				targetStmt = stmtNormal
			}
			if statusInt >= 200 && statusInt < 300 {
				batchSuccess++
			} else if statusInt >= 400 && statusInt < 500 {
				batchClientErr++
			} else if statusInt >= 500 {
				batchServerErr++
			}
		}

		rowUUID := uuid.NewString()
		_, err := targetStmt.ExecContext(ctx,
			rowUUID,
			logData.Timestamp, logData.Host, logData.ServerIP, logData.ClientIP, countryIso,
			logData.Xff, logData.Domain, logData.URL, logData.Referer, logData.Args,
			upstreamTimeStr, responseTimeStr, logData.RequestMethod, statusStr, sizeStr,
			logData.RequestBody, requestLengthParam, logData.Protocol, logData.UpstreamHost,
			logData.FileDir, logData.UserAgent, trafficType,
		)
		if err != nil {
			log.Printf("Failed to record IP data %s: %v", logData.ClientIP, err)
			continue
		}
		deletedDocIDs = append(deletedDocIDs, hit.ID)
	}

	todayStr := now.Format("2006-01-02")
	if len(batchAttackSummary) > 0 {
		stmtSummary, err := tx.PrepareContext(ctx, `
			INSERT INTO "ict_nginx_atc_sum" (
				id, date, client_ip, traffic_type, target_domain, total_hits, last_seen
			) VALUES ($1, $2, $3, $4, $5, $6, NOW())
			ON CONFLICT (date, client_ip, traffic_type, target_domain)
			DO UPDATE SET total_hits = "ict_nginx_atc_sum".total_hits + EXCLUDED.total_hits, last_seen = NOW()
		`)
		if err == nil {
			for k, hits := range batchAttackSummary {
				sumUUID := uuid.NewString()
				_, _ = stmtSummary.ExecContext(ctx,
					sumUUID, todayStr, k.ClientIP, k.TrafficType, k.Domain, hits,
				)
			}
			stmtSummary.Close()
		}
	}

	if batchTotal > 0 {
		avgTime := totalResponseTime / float64(batchTotal)
		slaUUID := uuid.NewString()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO "ict_nginx_sla" (id, date, total_requests, successful_requests, client_errors, server_errors, attack_requests, avg_response_time, sla_percentage)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0.00)
			ON CONFLICT (date) DO UPDATE SET
				total_requests = "ict_nginx_sla".total_requests + EXCLUDED.total_requests,
				successful_requests = "ict_nginx_sla".successful_requests + EXCLUDED.successful_requests,
				client_errors = "ict_nginx_sla".client_errors + EXCLUDED.client_errors,
				server_errors = "ict_nginx_sla".server_errors + EXCLUDED.server_errors,
				attack_requests = "ict_nginx_sla".attack_requests + EXCLUDED.attack_requests,
				avg_response_time = (("ict_nginx_sla".avg_response_time * "ict_nginx_sla".total_requests) + $9) / ("ict_nginx_sla".total_requests + EXCLUDED.total_requests)
		`,
			slaUUID, todayStr, batchTotal,
			batchSuccess, batchClientErr, batchServerErr, batchAttack,
			avgTime, totalResponseTime)
		if err != nil {
			log.Printf("Failed to process SLA: %v", err)
			return
		}
		_, _ = tx.ExecContext(ctx, `
			UPDATE "ict_nginx_sla"
			SET sla_percentage = CASE WHEN total_requests > 0 THEN (successful_requests::numeric / total_requests::numeric) * 100 ELSE 0 END
			WHERE date = $1
		`, todayStr)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to process transaction: %v", err)
		return
	}
	txOpened = false

	if len(deletedDocIDs) > 0 {
		idsJSON, _ := json.Marshal(deletedDocIDs)
		deleteQuery := fmt.Sprintf(`{"query":{"terms":{"_id":%s}}}`, idsJSON)
		resDel, err := es.DeleteByQuery([]string{indexName},
			strings.NewReader(deleteQuery),
			es.DeleteByQuery.WithContext(ctx),
			es.DeleteByQuery.WithConflicts("proceed"))
		drainBody(resDel.Body)
		if err != nil {
			log.Printf("Failed to bulk delete ES docs: %v", err)
		}
		log.Printf("Log %s deleted %d", indexName, len(deletedDocIDs))
	}
}

func cleanupOldIndices(ctx context.Context, es *elasticsearch.Client, now time.Time) {
	for i := 1; i <= 3; i++ {
		oldDate := now.AddDate(0, 0, -i)
		oldIndexName := fmt.Sprintf("logstash_%s", oldDate.Format("2006.01.02"))
		res, err := es.Indices.Delete([]string{oldIndexName}, es.Indices.Delete.WithContext(ctx))
		drainBody(res.Body)
		if err == nil && res != nil && res.StatusCode != http.StatusNotFound {
			log.Printf("Cleaned up old index: %s", oldIndexName)
		}
	}
}

func main() {
	cfg := LoadConfig()
	InitDB(cfg.PostgresConn)
	defer PgSQL.Close()

	CacheInstance.Refresh(PgSQL)

	esClient, err := initElasticClient()
	if err != nil {
		log.Fatalf("Failed to init Elasticsearch client: %v", err)
	}
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Failed to communicate with Elasticsearch cluster: %v", err)
	}
	drainBody(res.Body)

	StartSyncWorker(esClient, cfg.HttpFloodThreshold)
}
