package models

import (
	"database/sql"
	"time"
)

/*
tables

historical_2014_

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
        (ticker_id INT UNSIGNED NOT NULL,
         name VARCHAR(100) NOT NULL
        )
`

const createPriceTable = `
    CREATE TABLE IF NOT EXISTS 
        price
        (ticker_id INT UNSIGNED NOT NULL,
         price INT UNSIGNED NOT NULL,
         created_date DATETIME
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
        (user_id INT UNSIGNED NOT NULL,
         ticker_id INT UNSIGNED NOT NULL,
         shares BIGINT UNSIGNED NOT NULL,
         created_date DATETIME NOT NULL
        )
`

var allTables = []string{
	createTickerTable, createPriceTable,
	createArchivedPositionTable, createCurrentPositionTable,
}

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

func EnsureTables(db *sql.DB) error {
	for _, table := range allTables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}
