package models

import (
	"database/sql"
	"fmt"
	"go_trading/config"

	_ "github.com/mattn/go-sqlite3"
)

var DbConnection *sql.DB

func init() {
	/* DB接続 */
	DbConnection, _ = sql.Open("sqlite3", "stockdata.sql")
	/* テーブル作成(BTC_USD_1h0m0s, BTC_USD_1m0s, BTC_USD_1s) */
	for _, duration := range config.Durations {
		tableName := fmt.Sprintf("BTC_USD_%s", duration)
		cmd := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			time DATATIME PRIMARY KEY NOT NULL,
			open FLOAT,
			close FLOAT,
			high FLOAT,
			low FLOAT,
			volume FLOAT)`, tableName)
		DbConnection.Exec(cmd)
	}
}
