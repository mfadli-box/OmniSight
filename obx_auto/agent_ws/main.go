package main

import (
	"bufio"
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
	"path/filepath"
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

// ===== Shared =====

var (
	PgSQL    *sql.DB
	rotateMu sync.Mutex
	rotating bool
)

type Config struct {
	PostgresConn       string
	ElasticsearchURL   string
	HttpFloodThreshold int
	ArchiveDir         string
	NormalLogDays      int
	AttackLogDays      int
}

func LoadConfig() *Config {
	godotenv.Load()

	PG_Host := os.Getenv("PG_HOST")
	PG_Port := os.Getenv("PG_PORT")
	PG_User := os.Getenv("PG_USER")
	PG_Pass := os.Getenv("PG_PASS")
	PG_Data := os.Getenv("PG_DATA")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		PG_User, PG_Pass, PG_Host, PG_Port, PG_Data)

	thresh := 150
	if envThresh := os.Getenv("FT_HTTP"); envThresh != "" {
		if val, err := strconv.Atoi(envThresh); err == nil {
			thresh = val
		}
	}

	archiveDir := os.Getenv("RE_PATH")
	if archiveDir == "" {
		archiveDir = "./archive"
	}
	normalDays, _ := strconv.Atoi(os.Getenv("RE_NORMAL"))
	if normalDays <= 0 {
		normalDays = 7
	}
	attackDays, _ := strconv.Atoi(os.Getenv("RE_ATTACK"))
	if attackDays <= 0 {
		attackDays = 90
	}

	return &Config{
		PostgresConn:       dsn,
		ElasticsearchURL:   os.Getenv("ES_LINK"),
		HttpFloodThreshold: thresh,
		ArchiveDir:         archiveDir,
		NormalLogDays:      normalDays,
		AttackLogDays:      attackDays,
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

func drainBody(body io.ReadCloser) {
	if body == nil {
		return
	}
	defer body.Close()
	io.Copy(io.Discard, body)
}

// ===== Nginx Log Sync Agent =====

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
	CacheInstance   = &Cache{}
	TrackerInstance = &IPTracker{count: make(map[string]int)}
)

const cacheRefreshInterval = 5 * time.Minute

func (c *Cache) Refresh(db *sql.DB) {
	c.mu.Lock()
	defer c.mu.Unlock()

	exactIPs := make(map[string]bool)
	var cidrNets []*net.IPNet

	rows, err := db.Query(`SELECT ip_or_cidr FROM "ict_ip_whitelist"`)
	if err == nil && rows != nil {
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

	banned := make(map[string]bool)
	rowsBan, err := db.Query(`SELECT ip FROM "ict_ip_blacklist" WHERE expires_at IS NULL OR expires_at > NOW()`)
	if err == nil && rowsBan != nil {
		defer rowsBan.Close()
		for rowsBan.Next() {
			var ip string
			if err := rowsBan.Scan(&ip); err == nil {
				banned[ip] = true
			}
		}
	}

	var bypassRules []BypassRule
	rowsBypass, err := db.Query(`SELECT domain, url_path, COALESCE(args_pattern, '') FROM "ict_waf_bypass_rule"`)
	if err == nil && rowsBypass != nil {
		defer rowsBypass.Close()
		for rowsBypass.Next() {
			var r BypassRule
			if err := rowsBypass.Scan(&r.Domain, &r.URLPath, &r.ArgsPattern); err == nil {
				bypassRules = append(bypassRules, r)
			}
		}
	}

	c.exactIPs = exactIPs
	c.cidrNets = cidrNets
	c.banned = banned
	c.bypassRules = bypassRules
	c.loadedAt = time.Now()
	log.Printf("[nginx] Cache refreshed: %d whitelist IPs, %d CIDR, %d banned, %d bypass rules",
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
		domainMatch := r.Domain == "*" || r.Domain == domain
		urlMatch := false
		if strings.Contains(r.URLPath, "*") {
			pattern := strings.ReplaceAll(r.URLPath, "*", "%")
			urlMatch = matchLike(pattern, urlPath)
		} else {
			urlMatch = r.URLPath == urlPath
		}
		if domainMatch && urlMatch {
			if r.ArgsPattern == "" || args == "" || strings.Contains(args, r.ArgsPattern) {
				return true
			}
		}
	}
	return false
}

func matchLike(pattern, text string) bool {
	parts := strings.Split(pattern, "%")
	for _, part := range parts {
		if part != "" && !strings.Contains(text, part) {
			return false
		}
	}
	return true
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
			INSERT INTO "ict_ip_blacklist" (id, ip, threat_score, reason, expires_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (ip) DO UPDATE
			SET threat_score = EXCLUDED.threat_score, expires_at = EXCLUDED.expires_at
		`
		_, err := db.Exec(query, uuid.NewString(), ip, finalScore, trafficType, expiryTime)
		if err != nil {
			log.Printf("[nginx] Failed to save IP block %s: %v", ip, err)
		}
		CacheInstance.AddBanned(ip)
		return finalScore, true
	}
	return finalScore, false
}

func initElasticClient(addr string) (*elasticsearch.Client, error) {
	if addr == "" {
		addr = "http://elasticsearch:9200"
	}
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{addr},
		Transport: transport,
	})
}

func restartSelf() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalf("[nginx] Failed to find binary path: %v", err)
	}
	cmd := exec.Command(self, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		log.Fatalf("[nginx] Failed to restart: %v", err)
	}
	os.Exit(0)
}

func cleanupOldIndices(ctx context.Context, es *elasticsearch.Client, now time.Time) {
	for i := 1; i <= 3; i++ {
		oldDate := now.AddDate(0, 0, -i)
		oldIndexName := fmt.Sprintf("logstash_%s", oldDate.Format("2006.01.02"))
		res, err := es.Indices.Delete([]string{oldIndexName}, es.Indices.Delete.WithContext(ctx))
		drainBody(res.Body)
		if err == nil && res != nil && res.StatusCode != http.StatusNotFound {
			log.Printf("[nginx] Cleaned up old index: %s", oldIndexName)
		}
	}
}

func syncAndAnalyze(ctx context.Context, es *elasticsearch.Client, indexName string, now time.Time, threshold int) {
	cleanJsonQuery := `{"query":{"bool":{"must":{"match":{"error.message":"Error decoding JSON: invalid character 'x' in string escape code"}}}}}`
	res, _ := es.DeleteByQuery([]string{indexName}, strings.NewReader(cleanJsonQuery),
		es.DeleteByQuery.WithContext(ctx), es.DeleteByQuery.WithConflicts("proceed"))
	drainBody(res.Body)

	cleanStatus0 := `{"query":{"match":{"status":"0"}}}`
	res, _ = es.DeleteByQuery([]string{indexName}, strings.NewReader(cleanStatus0),
		es.DeleteByQuery.WithContext(ctx), es.DeleteByQuery.WithConflicts("proceed"))
	drainBody(res.Body)

	cleanIpQuery := `{"query":{"bool":{"should":[{"term":{"client_ip":""}},{"bool":{"must_not":{"exists":{"field":"client_ip"}}}}],"minimum_should_match":1}}}`
	res, _ = es.DeleteByQuery([]string{indexName}, strings.NewReader(cleanIpQuery),
		es.DeleteByQuery.WithContext(ctx), es.DeleteByQuery.WithConflicts("proceed"))
	drainBody(res.Body)

	queryFindDates := `
		SELECT DISTINCT date_str FROM (
			SELECT DISTINCT TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_str
			FROM   "ict_nginx_log"
			WHERE  client_ip = '' AND timestamp >= NOW() - INTERVAL '7 days'
			UNION
			SELECT DISTINCT TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_str
			FROM   "ict_nginx_app"
			WHERE  client_ip = '' AND timestamp >= NOW() - INTERVAL '7 days'
			UNION
			SELECT DISTINCT TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_str
			FROM   "ict_nginx_atc"
			WHERE  client_ip = '' AND timestamp >= NOW() - INTERVAL '7 days'
		) t WHERE date_str IS NOT NULL
	`
	rowsDates, errDates := PgSQL.QueryContext(ctx, queryFindDates)
	if errDates == nil && rowsDates != nil {
		var affectedDates []string
		for rowsDates.Next() {
			var d string
			if err := rowsDates.Scan(&d); err == nil {
				affectedDates = append(affectedDates, d)
			}
		}
		rowsDates.Close()
		if len(affectedDates) > 0 {
			log.Printf("[nginx] Cleaning client_ip for %d dates...", len(affectedDates))
			txClean, errTxClean := PgSQL.BeginTx(ctx, nil)
			if errTxClean == nil {
				cleanSuccess := true
				for _, dateStr := range affectedDates {
					_, err1 := txClean.ExecContext(ctx, `DELETE FROM "ict_nginx_log" WHERE client_ip = '' AND TO_CHAR(timestamp, 'YYYY-MM-DD') = $1`, dateStr)
					_, err2 := txClean.ExecContext(ctx, `DELETE FROM "ict_nginx_app" WHERE client_ip = '' AND TO_CHAR(timestamp, 'YYYY-MM-DD') = $1`, dateStr)
					_, err3 := txClean.ExecContext(ctx, `DELETE FROM "ict_nginx_atc" WHERE client_ip = '' AND TO_CHAR(timestamp, 'YYYY-MM-DD') = $1`, dateStr)
					_, err4 := txClean.ExecContext(ctx, `DELETE FROM "ict_nginx_atc_sum" WHERE client_ip = '' AND date = $1`, dateStr)
					if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
						continue
					}
					queryUpdateSLA := `
						WITH metrics AS (
							SELECT
								COUNT(*) AS total,
								COUNT(CASE WHEN (status::int >= 200 AND status::int < 300) AND (traffic_type = 'NORMAL' OR traffic_type = 'WHITELISTED_TRAFFIC') THEN 1 END) AS success,
								COUNT(CASE WHEN status::int >= 400 AND status::int < 500 AND status::int != 444 THEN 1 END) AS client_err,
								COUNT(CASE WHEN status::int >= 500 THEN 1 END) AS server_err,
								COALESCE(AVG(NULLIF(responsetime, '')::numeric), 0.0000) AS avg_time
							FROM (
								SELECT timestamp, status, traffic_type, responsetime FROM "ict_nginx_log" WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
								UNION ALL
								SELECT timestamp, status, traffic_type, responsetime FROM "ict_nginx_app" WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1 AND traffic_type != 'BLOCK_444'
								UNION ALL
								SELECT timestamp, status, traffic_type, responsetime FROM "ict_nginx_atc" WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
							) combined
						),
						attack_metrics AS (
							SELECT COUNT(*) AS total_attacks FROM "ict_nginx_atc"
							WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
							  AND traffic_type NOT LIKE 'WH_%%'
							  AND traffic_type NOT LIKE 'BR_%%'
							  AND NOT EXISTS(SELECT 1 FROM "ict_ip_blacklist" b WHERE b.ip = "ict_nginx_atc".client_ip AND b.is_permanent = true)
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
					log.Printf("[nginx] SLA updated for cleaned dates.")
				} else {
					txClean.Rollback()
				}
			}
		}
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
		return
	}
	var searchResult EsSearchResponse
	if err := json.NewDecoder(resSearch.Body).Decode(&searchResult); err != nil {
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
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`)
	if err != nil {
		return
	}
	defer stmtNormal.Close()

	stmtAttack, err := tx.PrepareContext(ctx, `
		INSERT INTO "ict_nginx_atc" (
			id, timestamp, host, server_ip, client_ip, country_iso, xff,
			domain, url, referer, args, upstreamtime, responsetime,
			request_method, status, size, request_body, request_length,
			protocol, upstreamhost, file_dir, http_user_agent, traffic_type
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`)
	if err != nil {
		return
	}
	defer stmtAttack.Close()

	stmtApp, err := tx.PrepareContext(ctx, `
		INSERT INTO "ict_nginx_app" (
			id, timestamp, host, server_ip, client_ip, country_iso, xff,
			domain, url, referer, args, upstreamtime, responsetime,
			request_method, status, size, request_body, request_length,
			protocol, upstreamhost, file_dir, http_user_agent, traffic_type
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23)`)
	if err != nil {
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
			trafficType = ClassifyTraffic(logData.URL, logData.Args, logData.RequestBody, logData.UserAgent)
		}
		batchTotal++

		isClientWhitelisted := CacheInstance.IsWhitelisted(logData.ClientIP)
		isRuleBypassed := CacheInstance.IsBypassRule(logData.Domain, logData.URL, logData.Args)
		isAttack := trafficType != "NORMAL"

		var targetStmt *sql.Stmt
		// Status 444 = nginx mitigation block, always insert to ict_nginx_app
		if statusInt == 444 {
			targetStmt = stmtApp
			trafficType = "BLOCK_444"
		} else if isAttack {
			if isClientWhitelisted {
				targetStmt = stmtApp
				trafficType = "WH_" + trafficType
			} else if isRuleBypassed {
				targetStmt = stmtApp
				trafficType = "BR_" + trafficType
			} else {
				targetStmt = stmtAttack
				batchAttack++
				key := AttackKey{ClientIP: logData.ClientIP, TrafficType: trafficType, Domain: logData.Domain}
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
		}

		if statusInt >= 200 && statusInt < 300 {
			batchSuccess++
		} else if statusInt >= 400 && statusInt < 500 {
			batchClientErr++
		} else if statusInt >= 500 {
			batchServerErr++
		}

		rowUUID := uuid.NewString()
		_, err := targetStmt.ExecContext(ctx,
			rowUUID, logData.Timestamp, logData.Host, logData.ServerIP, logData.ClientIP, countryIso,
			logData.Xff, logData.Domain, logData.URL, logData.Referer, logData.Args,
			upstreamTimeStr, responseTimeStr, logData.RequestMethod, statusStr, sizeStr,
			logData.RequestBody, requestLengthParam, logData.Protocol, logData.UpstreamHost,
			logData.FileDir, logData.UserAgent, trafficType,
		)
		if err != nil {
			continue
		}
		deletedDocIDs = append(deletedDocIDs, hit.ID)
	}

	todayStr := now.Format("2006-01-02")
	if len(batchAttackSummary) > 0 {
		stmtSummary, err := tx.PrepareContext(ctx, `
			INSERT INTO "ict_nginx_atc_sum" (id, date, client_ip, traffic_type, target_domain, total_hits, last_seen)
			VALUES ($1, $2, $3, $4, $5, $6, NOW())
			ON CONFLICT (date, client_ip, traffic_type, target_domain)
			DO UPDATE SET total_hits = "ict_nginx_atc_sum".total_hits + EXCLUDED.total_hits, last_seen = NOW()
		`)
		if err == nil {
			for k, hits := range batchAttackSummary {
				_, _ = stmtSummary.ExecContext(ctx, uuid.NewString(), todayStr, k.ClientIP, k.TrafficType, k.Domain, hits)
			}
			stmtSummary.Close()
		}
	}

	if batchTotal > 0 {
		avgTime := totalResponseTime / float64(batchTotal)
		_, _ = tx.ExecContext(ctx, `
			INSERT INTO "ict_nginx_sla" (id, date, total_requests, successful_requests, client_errors, server_errors, attack_requests, avg_response_time, sla_percentage)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0.00)
			ON CONFLICT (date) DO UPDATE SET
				total_requests = "ict_nginx_sla".total_requests + EXCLUDED.total_requests,
				successful_requests = "ict_nginx_sla".successful_requests + EXCLUDED.successful_requests,
				client_errors = "ict_nginx_sla".client_errors + EXCLUDED.client_errors,
				server_errors = "ict_nginx_sla".server_errors + EXCLUDED.server_errors,
				attack_requests = "ict_nginx_sla".attack_requests + EXCLUDED.attack_requests,
				avg_response_time = (("ict_nginx_sla".avg_response_time * "ict_nginx_sla".total_requests) + $9) / ("ict_nginx_sla".total_requests + EXCLUDED.total_requests)
		`, uuid.NewString(), todayStr, batchTotal, batchSuccess, batchClientErr, batchServerErr, batchAttack, avgTime, totalResponseTime)
		_, _ = tx.ExecContext(ctx, `
			WITH metrics AS (
				SELECT
					TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_key,
					COUNT(*) AS total,
					COUNT(CASE WHEN status::int >= 200 AND status::int < 300 THEN 1 END) AS success,
					COUNT(CASE WHEN status::int >= 400 AND status::int < 500 AND status::int != 444 THEN 1 END) AS client_err,
					COUNT(CASE WHEN status::int >= 500 THEN 1 END) AS server_err,
					COALESCE(AVG(NULLIF(responsetime, '')::numeric), 0.0000) AS avg_time
				FROM (
					SELECT timestamp, status, responsetime FROM "ict_nginx_log" WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
					UNION ALL
					SELECT timestamp, status, responsetime FROM "ict_nginx_app" WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1 AND traffic_type != 'BLOCK_444'
					UNION ALL
					SELECT timestamp, status, responsetime FROM "ict_nginx_atc" WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
				) ts
				GROUP BY TO_CHAR(timestamp, 'YYYY-MM-DD')
			),
			attack_metrics AS (
				SELECT
					TO_CHAR(timestamp, 'YYYY-MM-DD') AS date_key,
					COUNT(*) AS total_attacks
				FROM "ict_nginx_atc"
				WHERE TO_CHAR(timestamp, 'YYYY-MM-DD') = $1
				  AND traffic_type NOT LIKE 'WH_%%'
				  AND traffic_type NOT LIKE 'BR_%%'
				  AND NOT EXISTS(SELECT 1 FROM "ict_ip_blacklist" b WHERE b.ip = "ict_nginx_atc".client_ip AND b.is_permanent = true)
				GROUP BY TO_CHAR(timestamp, 'YYYY-MM-DD')
			)
			INSERT INTO "ict_nginx_sla" (id, date, total_requests, successful_requests, client_errors, server_errors, attack_requests, avg_response_time, sla_percentage)
			SELECT gen_random_uuid()::text, m.date_key, m.total, m.success, m.client_err, m.server_err, COALESCE(a.total_attacks, 0), m.avg_time,
				CASE WHEN m.total > 0 THEN ROUND(((m.total - m.server_err - COALESCE(a.total_attacks, 0))::numeric / m.total) * 100, 2) ELSE 0 END
			FROM metrics m
			LEFT JOIN attack_metrics a ON m.date_key = a.date_key
			ON CONFLICT (date) DO UPDATE SET
				total_requests = EXCLUDED.total_requests,
				successful_requests = EXCLUDED.successful_requests,
				client_errors = EXCLUDED.client_errors,
				server_errors = EXCLUDED.server_errors,
				attack_requests = EXCLUDED.attack_requests,
				avg_response_time = EXCLUDED.avg_response_time,
				sla_percentage = EXCLUDED.sla_percentage
		`, todayStr)
	}

	if err := tx.Commit(); err != nil {
		return
	}
	txOpened = false

	if len(deletedDocIDs) > 0 {
		idsJSON, _ := json.Marshal(deletedDocIDs)
		deleteQuery := fmt.Sprintf(`{"query":{"terms":{"_id":%s}}}`, idsJSON)
		resDel, _ := es.DeleteByQuery([]string{indexName},
			strings.NewReader(deleteQuery),
			es.DeleteByQuery.WithContext(ctx),
			es.DeleteByQuery.WithConflicts("proceed"))
		drainBody(resDel.Body)
	}
}

func runNginxSync(es *elasticsearch.Client, threshold int) {
	rotateMu.Lock()
	if rotating {
		rotateMu.Unlock()
		return
	}
	rotateMu.Unlock()

	ctx := context.Background()
	now := time.Now()
	eod := now.AddDate(0, 0, -1)
	beforeIndex := fmt.Sprintf("logstash_%s", eod.Format("2006.01.02"))
	currentIndex := fmt.Sprintf("logstash_%s", now.Format("2006.01.02"))
	syncAndAnalyze(ctx, es, beforeIndex, eod, threshold)
	syncAndAnalyze(ctx, es, currentIndex, now, threshold)
	cleanupOldIndices(ctx, es, eod)
}

// ===== Log Rotate Agent =====

var validTables = map[string]bool{
	"ict_nginx_log": true,
	"ict_nginx_app": true,
	"ict_nginx_atc": true,
}

const archiveBatchSize = 5000

func archiveTable(db *sql.DB, tableName string, archiveDir string, retentionDays int) int {
	if !validTables[tableName] {
		log.Printf("[rotate] Rejected archive of unknown table: %s", tableName)
		return 0
	}

	fileName := fmt.Sprintf("%s_archive_%s.log", tableName, time.Now().Format("2006-01-02"))
	filePath := filepath.Join(archiveDir, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[rotate] Failed to create archive file %s: %v", filePath, err)
		return 0
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 64*1024)
	defer writer.Flush()

	var totalCount int
	for {
		deleteQuery := fmt.Sprintf(`
			DELETE FROM %s WHERE ctid IN (
				SELECT ctid FROM %s
				WHERE  created_at < NOW() - $1 * INTERVAL '1 day' LIMIT %d
			) RETURNING row_to_json(%s)`,
			tableName, tableName, archiveBatchSize, tableName,
		)
		rows, err := db.Query(deleteQuery, retentionDays)
		if err != nil {
			log.Printf("[rotate] Failed to delete+return from %s: %v", tableName, err)
			return totalCount
		}
		var batchCount int
		for rows.Next() {
			var jsonRaw string
			if err := rows.Scan(&jsonRaw); err != nil {
				continue
			}
			writer.WriteString(jsonRaw)
			writer.WriteByte('\n')
			batchCount++
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return totalCount
		}
		rows.Close()
		totalCount += batchCount
		if batchCount < archiveBatchSize {
			break
		}
	}

	if totalCount > 0 {
		writer.Flush()
		log.Printf("[rotate] Archived %d rows from %s", totalCount, tableName)
	} else {
		os.Remove(filePath)
	}
	return totalCount
}

func vacuumTable(db *sql.DB, tableName string) {
	log.Printf("[rotate] Running VACUUM ANALYZE on %s...", tableName)
	_, err := db.Exec(fmt.Sprintf("VACUUM ANALYZE %s", tableName))
	if err != nil {
		log.Printf("[rotate] VACUUM failed on %s: %v", tableName, err)
	} else {
		log.Printf("[rotate] VACUUM ANALYZE completed on %s", tableName)
	}
}

func runLogRotate(cfg *Config) {
	rotateMu.Lock()
	rotating = true
	rotateMu.Unlock()

	defer func() {
		rotateMu.Lock()
		rotating = false
		rotateMu.Unlock()
	}()

	log.Println("[rotate] Starting log archiving and retention process...")
	if err := os.MkdirAll(cfg.ArchiveDir, 0755); err != nil {
		log.Fatalf("[rotate] Failed to create archive directory: %v", err)
	}

	archiveTable(PgSQL, "ict_nginx_log", cfg.ArchiveDir, cfg.NormalLogDays)
	archiveTable(PgSQL, "ict_nginx_app", cfg.ArchiveDir, cfg.AttackLogDays)
	archiveTable(PgSQL, "ict_nginx_atc", cfg.ArchiveDir, cfg.AttackLogDays)

	vacuumTable(PgSQL, "ict_nginx_log")
	vacuumTable(PgSQL, "ict_nginx_app")
	vacuumTable(PgSQL, "ict_nginx_atc")
	vacuumTable(PgSQL, "ict_nginx_atc_sum")
	vacuumTable(PgSQL, "ict_nginx_sla")

	log.Println("[rotate] Log archiving and retention process completed.")
}

func nextRunTime(hour, minute int) time.Duration {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.UTC)
	if !target.After(now) {
		target = target.Add(24 * time.Hour)
	}
	return target.Sub(now)
}

// ===== Main =====

func main() {
	cfg := LoadConfig()
	InitDB(cfg.PostgresConn)
	defer PgSQL.Close()

	CacheInstance.Refresh(PgSQL)

	esClient, err := initElasticClient(cfg.ElasticsearchURL)
	if err != nil {
		log.Fatalf("Failed to init Elasticsearch client: %v", err)
	}
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Failed to communicate with Elasticsearch cluster: %v", err)
	}
	drainBody(res.Body)

	// Nginx sync: every 1 minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		cacheTicker := time.NewTicker(cacheRefreshInterval)
		for {
			select {
			case <-ticker.C:
				runNginxSync(esClient, cfg.HttpFloodThreshold)
			case <-cacheTicker.C:
				CacheInstance.Refresh(PgSQL)
			}
		}
	}()

	// Log rotate: daily at 02:00 UTC
	go func() {
		delay := nextRunTime(2, 0)
		log.Printf("[rotate] Next rotation in %s", delay)
		time.Sleep(delay)
		runLogRotate(cfg)
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			runLogRotate(cfg)
		}
	}()

	log.Println("agent_ws started. Nginx sync: 1min, Log rotate: 02:00 UTC daily.")
	select {}
}
