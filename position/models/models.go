package models

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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

	// starting price of all tickers in mils
	startPrice = 10000

	//
	priceMin = 1000
	priceMax = 20000

	// minimum time between new stock prices
	priceChangeSecMin = 3000
	priceChangeSecMax = 60000

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

	getTickers = `SELECT ticker_id FROM ticker`

	insertPrice = `
    INSERT INTO price (ticker_id, price, created_date) 
		SELECT t.ticker_id, ?, ? FROM ticker t
			WHERE t.symbol = ?
	ON DUPLICATE KEY UPDATE price.ticker_id = price.ticker_id
    `

	//  alter table price add partition if not exists (partition p201002 values less than ('2010-02-01'));
	addPartition = "ALTER TABLE `%s` ADD PARTITION IF NOT EXISTS (PARTITION `%s` VALUES LESS THAN ('%s'))"

	createPriceTable = `
    CREATE TABLE IF NOT EXISTS 
        price
        (ticker_id INT UNSIGNED NOT NULL,
         price INT UNSIGNED NOT NULL,
         created_date DATETIME,
         UNIQUE KEY ticker_date (ticker_id, created_date)
        )
        PARTITION BY RANGE COLUMNS(created_date) (
            PARTITION p201002 VALUES LESS THAN ('2010-02-01')
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

var startDate = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC)

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
			log.Println(table)
			log.Println(err)
			return err
		}
	}

	return nil
}

// Load ticker table and set initial price
func LoadTickers(db *sql.DB) error {
	for _, ticker := range strings.Fields(tickerSymbols) {
		_, err := db.Exec(insertTicker, ticker)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = db.Exec(insertPrice, startPrice, startDate, ticker)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func nextDate(date time.Time) time.Time {
	v := rand.Int31n(priceChangeSecMax-priceChangeSecMin) + priceChangeSecMin
	return date.Add(time.Duration(v) * time.Second)
}

// Return a time in the future that is the beginning of the next month.
func partitionDate(date time.Time) time.Time {
	var newDate time.Time
	if date.Month() == time.December {
		newDate = time.Date(date.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	} else {
		newDate = time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	}

	return newDate
}

// Load random prices from date range
func LoadPrices(db *sql.DB) error {
	date := startDate
	for i := 0; i < 100; i++ {
		date = nextDate(date)
		pDate := partitionDate(date)

		partitionCommand := fmt.Sprintf(addPartition, "price",
			pDate.Format("p200601"), pDate.Format("2006-01-02"),
		)

		_, err := db.Exec(partitionCommand)
		if err != nil {
			log.Println(partitionCommand)
			log.Println(date)
			log.Println(pDate)
			log.Println(err)
			return err
		}

		for _, ticker := range strings.Fields(tickerSymbols) {
			newPrice := rand.Int31n(priceMax-priceMin) + priceMin
			_, err := db.Exec(insertPrice, newPrice, date, ticker)
			if err != nil {
				log.Println(err)
				return err
			}

		}
	}

	return nil
}
