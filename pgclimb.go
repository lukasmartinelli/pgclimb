package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/lukasmartinelli/pgclimb/formats"
	"github.com/lukasmartinelli/pgclimb/pg"
)

func changeHelpTemplateArgs(args string) {
	cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", args, -1)
}

func isSqlFile(arg string) bool {
	hasSelect := strings.HasPrefix(strings.ToLower(arg), "select")
	hasSqlExtension := strings.HasSuffix(arg, ".sql")
	return hasSqlExtension && !hasSelect
}

func parseQuery(c *cli.Context, command string) string {
	arg := c.Args().First()
	if arg == "" {
		cli.ShowCommandHelp(c, command)
		os.Exit(1)
	}

	if isSqlFile(arg) {
		filename := arg
		query, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalln(err)
		}
		return string(query)
	} else {
		return arg
	}
}

func exitOnError(err error) {
	log.SetFlags(0)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "pgclimb"
	app.Version = "0.1"
	app.Usage = "Export data from PostgreSQL into different data formats"
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
				connStr := pg.ParseConnStr(c)
				err := formats.Export(query, connStr, formats.NewJSONArrayFormat())
				exitOnError(err)
			},
		},
		{
			Name:  "json",
			Usage: "Export JSON document",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "json")
				connStr := pg.ParseConnStr(c)
				err := formats.Export(query, connStr, formats.NewJSONArrayFormat())
				exitOnError(err)
			},
		},
		{
			Name:  "csv",
			Usage: "Export CSV",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "csv")
				connStr := pg.ParseConnStr(c)
				err := formats.Export(query, connStr, formats.NewCsvFormat(';'))
				exitOnError(err)
			},
		},
		{
			Name:  "tsv",
			Usage: "Export TSV",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "tsv")
				connStr := pg.ParseConnStr(c)
				err := formats.Export(query, connStr, formats.NewCsvFormat('\t'))
				exitOnError(err)
			},
		},
		{
			Name:  "xml",
			Usage: "Export XML",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "xml")
				connStr := pg.ParseConnStr(c)
				err := formats.Export(query, connStr, formats.NewXMLFormat())
				exitOnError(err)
			},
		},
		{
			Name:  "xlsx",
			Usage: "Export XLSX spreadsheets",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<query>")
				query := parseQuery(c, "xlsx")
				connStr := pg.ParseConnStr(c)
				err := formats.Export(query, connStr, formats.NewXlsxFormat())
				exitOnError(err)
			},
		},
	}

	app.Run(os.Args)
}
