package models

import (
	"database/sql"
	"strings"
	"time"
)

/*
tables

current

operations
    AddPosition(userId, tickerId, shares, date)
    Archive()
    GetPrices()
    Execute()
    GetPresentBalances()
    GetHistoricalBalances(startDate, endDate)

*/

const (
	tickerSymbols = `SPY VXX EWZ EEM QQQ IWM XLF TVIX XLE XIV TZA EWJ GDX UVXY EFA FXI XOP XLU VWO XLV SDS OIH XLP RSX XLI XLK TNA GDXJ XLB DGAZ IYR UGAZ IVV EWT HYG DIA TLT USO EWG JNUG SSO SPXU SQQQ NUGT JNK XLY XHB QID USMV UNG VGK SH ERY VEA AMLP IAU FAZ DXJ GLD TQQQ ITB EWY VNQ SLV EPI VTI EWP XRT SPXS DUST EDC SMH SVXY IWD BND BSV EZU AAXJ TWM QLD IEMG EWH JDST TBT BKLN IWF EWU LQD DXD VIXY RWM EWA MDY SPXL DBC REM FAS KRE IBB SCO`

	createTickerTable = `
    CREATE TABLE IF NOT EXISTS
        ticker
        (ticker_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
         symbol VARCHAR(4) NOT NULL,
         UNIQUE KEY symbol_key (symbol)
        )
    `

	insertTicker = `
    INSERT INTO ticker (symbol) VALUES(?)
    ON DUPLICATE KEY UPDATE ticker_id = ticker_id
    `

	createPriceTable = `
    CREATE TABLE IF NOT EXISTS 
        price
        (ticker_id INT UNSIGNED NOT NULL,
         price INT UNSIGNED NOT NULL,
         created_date DATETIME,
         UNIQUE KEY ticker_date (ticker_id, created_date)
        )
        PARTITION BY RANGE COLUMNS(created_date) (
            PARTITION p201001 VALUES LESS THAN ('2010-01-01'),
            PARTITION pMax VALUES LESS THAN MAXVALUE
        )
    `

	// partitioned by date, indexed by user id?
	createArchivedPositionTable = `
    CREATE TABLE IF NOT EXISTS 
        archived_position
        (user_id INT UNSIGNED NOT NULL,
         ticker_id INT UNSIGNED NOT NULL,
         shares BIGINT UNSIGNED NOT NULL,
         created_date DATETIME NOT NULL,
         price INT UNSIGNED 
        )
    `

	// partitioned and indexed by user id
	createCurrentPositionTable = ` 
    CREATE TABLE IF NOT EXISTS 
        position
        (user_id INT UNSIGNED NOT NULL PRIMARY KEY,
         ticker_id INT UNSIGNED NOT NULL,
         shares BIGINT UNSIGNED NOT NULL,
         created_date DATETIME NOT NULL
        )
    `
)

var allTables = []string{
	createTickerTable, createPriceTable,
	createArchivedPositionTable, createCurrentPositionTable,
}

var startDate = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.Local)

type Position struct {
	UserId   int
	TickerId int

	// Shares
	Shares uint64

	CreatedDate *time.Time

	Price sql.NullInt64
}

type Price struct {
	TickerId int
	Price    int
	Date     *time.Time
}

type Ticker struct {
	TickerId int
	Symbol   string
}

func EnsureTables(db *sql.DB) error {
	for _, table := range allTables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadTickers(db *sql.DB) error {
	for _, ticker := range strings.Fields(tickerSymbols) {
		_, err := db.Exec(insertTicker, ticker)
		if err != nil {
			return err
		}
	}

	return nil
}
