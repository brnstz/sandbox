package models

import (
	"database/sql"
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

const createTickerTable = `
    CREATE TABLE IF NOT EXISTS
        ticker
        (ticker_id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
         symbol VARCHAR(4) NOT NULL
        )
`

const createPriceTable = `
    CREATE TABLE IF NOT EXISTS 
        price
        (ticker_id INT UNSIGNED NOT NULL,
         price INT UNSIGNED NOT NULL,
         created_date DATETIME
        )
        UNIQUE KEY ticker_date (ticker_id, created_date)
        PARTITIONED BY RANGE COLUMNS(created_date) (
            PARTITION p201001 VALUES LESS THAN ('2010-01-01'),
            PARTITION pMax VALUES LESS THAN MAXVALUE
        )
`

// partitioned by date, indexed by user id?
const createArchivedPositionTable = `
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
const createCurrentPositionTable = ` 
    CREATE TABLE IF NOT EXISTS 
        position
        (user_id INT UNSIGNED NOT NULL PRIMARY KEY,
         ticker_id INT UNSIGNED NOT NULL,
         shares BIGINT UNSIGNED NOT NULL,
         created_date DATETIME NOT NULL
        )
`

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
