package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/brnstz/sandbox/position/models"
)

func main() {
	db, err := sql.Open("mysql", "bseitz:bseitz@localhost/brnstz")
	if err != nil {
		panic(err)
	}

	err = models.EnsureTables(db)
	if err != nil {
		panic(err)
	}
}
