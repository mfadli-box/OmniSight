package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var validTables = map[string]bool{
	"ict_nginx_log": true,
	"ict_nginx_app": true,
	"ict_nginx_atc": true,
}

const archiveBatchSize = 5000

type RetentionConfig struct {
	ArchiveDir    string
	NormalLogDays int
	AttackLogDays int
}

func LoadRetentionConfig() *RetentionConfig {
	dir := os.Getenv("RE_PATH")
	if dir == "" {
		dir = "./archive"
	}
	normalDays, _ := strconv.Atoi(os.Getenv("RE_NORMAL"))
	if normalDays <= 0 {
		normalDays = 7
	}
	attackDays, _ := strconv.Atoi(os.Getenv("RE_ATTACK"))
	if attackDays <= 0 {
		attackDays = 90
	}
	return &RetentionConfig{
		ArchiveDir:    dir,
		NormalLogDays: normalDays,
		AttackLogDays: attackDays,
	}
}

var PgSQL *sql.DB

func StartAutoArchiveAndRetentionWorker(db *sql.DB) {
	cfg := LoadRetentionConfig()
	if err := os.MkdirAll(cfg.ArchiveDir, 0755); err != nil {
		log.Fatalf("Failed to create archive directory: %v", err)
	}
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		log.Println("Starting log archiving and retention process...")
		archiveTable(db, "ict_nginx_log", cfg.ArchiveDir, cfg.NormalLogDays)
		archiveTable(db, "ict_nginx_app", cfg.ArchiveDir, cfg.AttackLogDays)
		archiveTable(db, "ict_nginx_atc", cfg.ArchiveDir, cfg.AttackLogDays)
		<-ticker.C
	}
}

func archiveTable(db *sql.DB, tableName string, archiveDir string, retentionDays int) {
	if !validTables[tableName] {
		log.Printf("Rejected archive of unknown table: %s", tableName)
		return
	}

	fileName := fmt.Sprintf("%s_archive_%s.log", tableName, time.Now().Format("2006-01-02"))
	filePath := filepath.Join(archiveDir, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to create archive file %s: %v", filePath, err)
		return
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
			log.Printf("Failed to delete+return from %s: %v", tableName, err)
			return
		}
		var batchCount int
		for rows.Next() {
			var jsonRaw string
			if err := rows.Scan(&jsonRaw); err != nil {
				log.Printf("Failed to scan row from %s: %v", tableName, err)
				continue
			}
			writer.WriteString(jsonRaw)
			writer.WriteByte('\n')
			batchCount++
		}
		if err := rows.Err(); err != nil {
			log.Printf("Row iteration error for %s: %v", tableName, err)
			rows.Close()
			return
		}
		rows.Close()
		totalCount += batchCount
		if batchCount < archiveBatchSize {
			break
		}
		log.Printf("Archived batch of %d rows from %s...", batchCount, tableName)
	}

	if totalCount > 0 {
		writer.Flush()
		log.Printf("Successfully archived %d rows from %s to %s", totalCount, tableName, filePath)
	} else {
		log.Printf("No old data in table %s needs to be archived today.", tableName)
		os.Remove(filePath)
	}
}

func main() {
	godotenv.Load()
	PG_Host := os.Getenv("PG_HOST")
	PG_Port := os.Getenv("PG_PORT")
	PG_User := os.Getenv("PG_USER")
	PG_Pass := os.Getenv("PG_PASS")
	PG_Data := os.Getenv("PG_DATA")
	IS_Pool := os.Getenv("IS_POOL")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		PG_Host, PG_Port, PG_User, PG_Pass, PG_Data)

	var err error
	PgSQL, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if IS_Pool == "true" {
		PgSQL.SetMaxOpenConns(100)
		PgSQL.SetMaxIdleConns(10)
	} else {
		PgSQL.SetMaxOpenConns(50)
		PgSQL.SetMaxIdleConns(25)
	}
	PgSQL.SetConnMaxLifetime(5 * time.Minute)

	if err = PgSQL.Ping(); err != nil {
		log.Fatalf("Database did not respond: %v", err)
	}

	StartAutoArchiveAndRetentionWorker(PgSQL)
}
