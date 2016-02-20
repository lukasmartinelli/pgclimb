package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/andrew-d/go-termutil"
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

func isTplFile(arg string) bool {
	return strings.HasSuffix(arg, ".tpl")
}

func parseTemplate(arg string) string {

	if isTplFile(arg) {
		filename := arg
		rawTemplate, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalln(err)
		}
		return string(rawTemplate)
	} else {
		return arg
	}
}

func exportFormat(c *cli.Context, format formats.DataFormat) error {
	connStr := pg.ParseConnStr(c)
	query, err := parseQuery(c)
	if err != nil {
		return err
	}

	return formats.Export(query, connStr, format)
}

func parseQuery(c *cli.Context) (string, error) {
	filename := c.String("file")
	if filename != "" {
		query, err := ioutil.ReadFile(filename)
		return string(query), err
	}

	command := c.String("command")
	if command != "" {
		return command, nil
	}

	if !termutil.Isatty(os.Stdin.Fd()) {
		query, err := ioutil.ReadAll(os.Stdin)
		return string(query), err
	}

	return "", errors.New("You need to specify a SQL query.")
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
			Name:   "dbname, d",
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
			Name:   "port, p",
			Value:  "5432",
			Usage:  "port",
			EnvVar: "DB_PORT",
		},
		cli.StringFlag{
			Name:   "username, U",
			Value:  "postgres",
			Usage:  "username",
			EnvVar: "DB_USER",
		},
		cli.BoolFlag{
			Name:  "ssl",
			Usage: "require ssl mode",
		},
		cli.StringFlag{
			Name:   "password, pass",
			Value:  "",
			Usage:  "password",
			EnvVar: "DB_PASS",
		},
		cli.StringFlag{
			Name:   "query, command, c",
			Value:  "",
			Usage:  "SQL query to execute",
			EnvVar: "DB_QUERY",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "SQL query filename",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "",
			Usage: "Output filename",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "template",
			Usage: "Export data with custom template",
			Action: func(c *cli.Context) {
				changeHelpTemplateArgs("<template>")

				templateArg := c.Args().First()
				if templateArg == "" {
					cli.ShowCommandHelp(c, "template")
					os.Exit(1)
				}
				rawTemplate := parseTemplate(templateArg)

				err := exportFormat(c, formats.NewTemplateFormat(rawTemplate))
				exitOnError(err)
			},
		},
		{
			Name:  "jsonlines",
			Usage: "Export newline-delimited JSON objects",
			Action: func(c *cli.Context) {
				err := exportFormat(c, formats.NewJSONLinesFormat())
				exitOnError(err)
			},
		},
		{
			Name:  "json",
			Usage: "Export JSON document",
			Action: func(c *cli.Context) {
				err := exportFormat(c, formats.NewJSONArrayFormat())
				exitOnError(err)
			},
		},
		{
			Name:  "csv",
			Usage: "Export CSV",
			Action: func(c *cli.Context) {
				err := exportFormat(c, formats.NewCsvFormat(';'))
				exitOnError(err)
			},
		},
		{
			Name:  "tsv",
			Usage: "Export TSV",
			Action: func(c *cli.Context) {
				err := exportFormat(c, formats.NewCsvFormat('\t'))
				exitOnError(err)
			},
		},
		{
			Name:  "xml",
			Usage: "Export XML",
			Action: func(c *cli.Context) {
				err := exportFormat(c, formats.NewXMLFormat())
				exitOnError(err)
			},
		},
		{
			Name:  "xlsx",
			Usage: "Export XLSX spreadsheets",
			Action: func(c *cli.Context) {
				err := exportFormat(c, formats.NewXlsxFormat())
				exitOnError(err)
			},
		},
	}

	app.Run(os.Args)
}
