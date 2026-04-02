package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "/var/lib/ascii-monitor/alerts.db"

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS alerts (
		hash      TEXT PRIMARY KEY,
		first_at  DATETIME NOT NULL,
		last_at   DATETIME NOT NULL,
		count     INTEGER  NOT NULL DEFAULT 1
	)`)
	return db, err
}

// isNewAlert 判断该错误是否为「新告警」（24小时内未出现过返回 true）
func isNewAlert(db *sql.DB, content string) (bool, error) {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))[:16]
	now := time.Now()
	cutoff := now.Add(-24 * time.Hour)

	var lastAt time.Time
	err := db.QueryRow("SELECT last_at FROM alerts WHERE hash = ?", hash).Scan(&lastAt)
	if err == sql.ErrNoRows {
		// 首次出现
		_, err = db.Exec("INSERT INTO alerts(hash,first_at,last_at,count) VALUES(?,?,?,1)", hash, now, now)
		return true, err
	}
	if err != nil {
		return false, err
	}

	// 更新 last_at 和 count
	_, _ = db.Exec("UPDATE alerts SET last_at=?, count=count+1 WHERE hash=?", now, hash)

	// 24小时内已告警过
	if lastAt.After(cutoff) {
		return false, nil
	}
	return true, nil
}
