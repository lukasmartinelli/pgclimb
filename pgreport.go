package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func exitOnError(err error) {
	log.SetFlags(0)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "pgreport"
	app.Version = "0.1"
	app.Usage = "Generate reports from PostgreSQL"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "dbname, db",
			Value:  "postgres",
			Usage:  "database",
			EnvVar: "DB_NAME",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			Usage:  "host name",
			EnvVar: "DB_HOST",
		},
		cli.StringFlag{
			Name:   "port",
			Value:  "5432",
			Usage:  "port",
			EnvVar: "DB_PORT",
		},
		cli.StringFlag{
			Name:   "username, user",
			Value:  "postgres",
			Usage:  "username",
			EnvVar: "DB_USER",
		},
		cli.BoolFlag{
			Name:  "ssl",
			Usage: "require ssl mode",
		},
		cli.StringFlag{
			Name:   "pass, pw",
			Value:  "",
			Usage:  "password",
			EnvVar: "DB_PASS",
		},
	}

	app.Action = func(c *cli.Context) {
		query := c.Args().First()
		if query == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		connStr := parseConnStr(c)
		err := export(query, connStr, NewCsvEncoder())
		exitOnError(err)
	}

	app.Run(os.Args)
}
