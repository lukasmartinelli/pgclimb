# pgclimb [![Build Status](https://travis-ci.org/lukasmartinelli/pgclimb.svg?branch=master)](https://travis-ci.org/lukasmartinelli/pgclimb) [![Go Report Card](https://goreportcard.com/badge/github.com/lukasmartinelli/pgclimb)](https://goreportcard.com/report/github.com/lukasmartinelli/pgclimb) ![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)

<img align="right" alt="Climbing elephant" src="logo.png" />

A PostgreSQL utility to export data into different data formats with
support for templates and an easier workflow than simply using `psql`.

Features:
- Export data to [JSON](#json-document), [JSON Lines](#json-lines), [CSV](#csv-and-tsv), [XLSX](#xlsx), [XML](#xml)
- [Templates](#templates) are supported for anything else (HTML, Markdown, Text)
- Advanced directory structures for generating e.g. readonly JSON APIs or multiple CSV files

Use Cases:
- `psql` alternative for getting data out of PostgreSQL
- Publish data sets
- Transform data for graphing it with spreadsheet software or JavaScript libraries
- Generate readonly JSON APIs
- Generate a web page

## Install

You can download a single binary for Linux, OSX or Windows.

**OSX**

```bash
wget -O pgclimb https://github.com/lukasmartinelli/pgfutter/releases/download/v0.3.2/pgfutter_darwin_amd64
chmod +x pgclimb

./pgclimb --help
```

**Linux**

```bash
wget -O pgclimb https://github.com/lukasmartinelli/pgfutter/releases/download/v0.3.2/pgfutter_linux_amd64
chmod +x pgclimb

./pgclimb --help
```

**Install from source**

```bash
go get github.com/lukasmartinelli/pgclimb
```

If you are using Windows or 32-bit architectures you need to [download the appropriate binary
yourself](https://github.com/lukasmartinelli/pgclimb/releases/latest).

## Supported Formats

The example queries operate on the open data [employee salaries of Montgomery County Maryland](https://data.montgomerycountymd.gov/Human-Resources/Employee-Salaries-2014/54rh-89p8).

## CSV and TSV

Exporting CSV and TSV files is very similar to using `psql` and the `COPY TO` statement. You can customize the delimiter which is `,` by default.

```bash
# Create a standard CSV file
pgclimb csv "SELECT * FROM employee_salaries"
# Create CSV file with custom delimiter and header row
pgclimb csv "SELECT full_name, position_title FROM employee_salaries" \
     --delimiter ";" --header
# Create TSV files
pgclimb tsv "SELECT position_title, COUNT(*) FROM employee_salaries GROUP BY position_title"
```

### JSON Document

Creating a single JSON document of a query is especially helpful if you
interface with other programs like providing data for JavaScript or creating
a readonly JSON API. You don't need to `json_agg` your objects, `pgclimb` will
automatically serialize the JSON for you - it does however supported nested JSON objects for more complicated queries.

```bash
# Query all salaries into JSON array
pgclimb json "SELECT * FROM employee_salaries"
# Render all employees of a position as nested JSON object
pgclimb json "SELECT s.position_title, json_agg(s) FROM employee_salaries s GROUP BY s.position_title"
```

### JSON Lines

[Newline delimited JSON](http://jsonlines.org/) is a good format to exchange
structured data in large quantities which does not fit well into the CSV format.

```bash
# Query all salaries into JSON array
pgclimb jsonlines "SELECT * FROM employee_salaries"
# Render all employees of a position as nested JSON object
pgclimb jsonlines "SELECT s.position_title, json_agg(s) FROM employee_salaries s GROUP BY s.position_title"
```

### XLSX

Excel files are useful for non programmers to directly work with the data
and create graphs and filters. You can also fill different data into different spreedsheets.

```bash
# Store all salaries in XLSX file
pgclimb xlsx "SELECT * FROM employee_salaries"
# Explicitly name sheet name
pgclimb xlsx "SELECT * FROM employee_salaries" --sheet "salaries"
```

### XML

You can output XML to process it with other programs or a XLST stylesheet.
If want more control over the XML output you can use the templating functionality
of `pgclimb` or build your own XML document with [XML functions in PostgreSQL](https://wiki.postgresql.org/wiki/XML_Support).

```bash
pgclimb xml "SELECT * FROM employee_salaries"
```

## Templates

This is the most advanced option and allows you to implement a lot of other formats and endpoints for free.
Because the template and query in this example are larger we fall back to using files instead of passing arguments.

Create a template `salaries.tpl`.

```html
<!DOCTYPE html>
<html>
    <head><title>Montgomery County MD Employees</title></head>
    <body>
        <h2>Employees</h2>
        <ul>
            {{range .}}
            <li>{{.full_name}}</li>
            {{end}}
        </ul>
    </body>
</html>
```

Create a query file `query.sql`

```sql
SELECT * FROM employee_salaries
```

And now run the template.

```
pgclimb template salaries.tpl query.sql
```

## Database Connection

Database connection details can be provided via environment variables
or as separate flags.

name        | default     | description
------------|-------------|------------------------------
`DB_NAME`   | `postgres`  | database name
`DB_HOST`   | `localhost` | host name
`DB_PORT`   | `5432`      | port
`DB_SCHEMA` | `import`    | schema to create tables for
`DB_USER`   | `postgres`  | database user
`DB_PASS`   |             | password (or empty if none)

## Advanced Use Cases

### Load SQL from File

If you have a long SQL statement to select your data you can read
the query from a file. Instead of passing a query to `pgclimb` you 
pass a filename ending with `.sql`.

```bash
# Store query in file
echo 'SELECT * FROM communities' > myquery.sql
# Execute query from file
pgclimb jsonlines myquery.sql
```

## Control filename via column data

Let's generate a `communities.json` files containing an overview of all
files and a file for each community containing the details.

Generate a single document.

```bash
pgclimb json "SELECT 'communities.json' AS filename, \\
          json_agg(t) AS document \\
          FROM (SELECT bfs_id, name FROM communities) AS t" --fname-field 'filename'
```

Generate multiple documents with the details.

```bash
pgclimb json "SELECT 'communities/' || bfs_id || '.json' AS filename, \
                 json_agg(c) AS document \
          FROM communities) AS c" --fname-field 'filename'
```

## Using JSON aggregation

This is not a `pgclimb` feature but shows you how to create more complex
JSON objects by using the [PostgreSQL JSON functions](http://www.postgresql.org/docs/9.5/static/functions-json.html).

Let's query communities and join an additional birth rate table.

```bash
pgclimb jsonlines "SELECT id, name, \\
    (SELECT array_to_json(array_agg(t)) FROM ( \\
            SELECT year, births FROM public.births \\
            WHERE community_id = c.id \\
            ORDER BY year ASC \\
        ) AS t \\
    ) AS births, \\
    FROM communities) AS c"
```

## Cross-compiling

We use [gox](https://github.com/mitchellh/gox) to create distributable
binaries for Windows, OSX and Linux.

```bash
docker run --rm -v "$(pwd)":/usr/src/pgclimb -w /usr/src/pgclimb tcnksm/gox:1.4.2-light
```
