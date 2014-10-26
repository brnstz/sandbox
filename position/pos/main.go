package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/brnstz/sandbox/position/models"
)

func main() {
	log.SetFlags(log.Llongfile)
	db, err := sql.Open("mysql", "bseitz:bseitz@/brnstz")
	if err != nil {
		panic(err)
	}

	err = models.EnsureTables(db)
	if err != nil {
		panic(err)
	}

	err = models.LoadTickers(db)
	if err != nil {
		panic(err)
	}

	err = models.LoadPrices(db)
	if err != nil {
		panic(err)
	}
}
