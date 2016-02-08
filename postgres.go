package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
)

//setup a database connection and create the import schema
func connect(connStr string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}

//parse sql connection string from cli flags
func parseConnStr(c *cli.Context) string {
	otherParams := "sslmode=disable connect_timeout=5"
	if c.GlobalBool("ssl") {
		otherParams = "sslmode=require connect_timeout=5"
	}
	return fmt.Sprintf("user=%s dbname=%s password='%s' host=%s port=%s %s",
		c.GlobalString("username"),
		c.GlobalString("dbname"),
		c.GlobalString("pass"),
		c.GlobalString("host"),
		c.GlobalString("port"),
		otherParams,
	)
}
