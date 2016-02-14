package main

import (
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

func changeHelpTemplateArgs(args string) {
	cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", args, -1)
}

func parseQuery(c *cli.Context, command string) string {
	query := c.Args().First()
	if query == "" {
		cli.ShowCommandHelp(c, command)
		os.Exit(1)
	}
	return query
}

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

	app.Commands = []cli.Command{
		{
			Name:  "jsonlines",
			Usage: "Export newline-delimited JSON objects",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "jsonlines")
				connStr := parseConnStr(c)
				err := export(query, connStr, NewJsonEncoder())
				exitOnError(err)
			},
		},
		{
			Name:  "csv",
			Usage: "Export CSV",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "csv")
				connStr := parseConnStr(c)
				err := export(query, connStr, NewCsvEncoder())
				exitOnError(err)
			},
		},
		{
			Name:  "xml",
			Usage: "Export XML",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "xml")
				connStr := parseConnStr(c)
				err := export(query, connStr, NewXmlEncoder())
				exitOnError(err)
			},
		},
	}

	app.Run(os.Args)
}
