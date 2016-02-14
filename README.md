# pgreport ![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)

<img align="right" alt="Climbing elephant" src="logo.png" />

A PostgreSQL utility to export data into different data formats with
support for templates and an easier workflow than simply using `psql`.

Features:
- Export data to JSON, JSON Lines, CSV, XLSX, XML
- Templates are supported for anything else (HTML, Markdown, Text)
- Advanced directory structures for generating e.g. readonly JSON APIs or multiple CSV files

Use Cases:
- `psql` alternative for getting data out of PostgreSQL
- Publish data sets
- Transform data for graphing it with spreadsheet software or JavaScript libraries
- Generate readonly JSON APIs
- Generate a web page 

## Examples

## Generate JSON Lines

[Newline delimited JSON](http://jsonlines.org/) is a good format to exchange
structured data in large quantities.

`pgreport` supports rendering JSON output for arbitrary queries. If you
want to export more complicated structured you can create JSON aggregation
in PostgreSQL and `pgreport` will handle it just fine.

Let's query communities join an additional birth rate table.

```bash
pgreport jsonlines "SELECT id, name, \\
    (SELECT array_to_json(array_agg(t)) FROM ( \\
            SELECT year, births FROM public.births \\
            WHERE community_id = c.id \\
            ORDER BY year ASC \\
        ) AS t \\
    ) AS births, \\
    FROM communities) AS c"
```

## Generate a readonly API

Let's generate a `communities.json` files containing an overview of all
files and a file for each community containing the details.

Generate a single document.

```bash
pgreport "SELECT 'communities.json' AS filename, \\
          json_agg(t) AS document \\
          FROM (SELECT bfs_id, name FROM communities) AS t"
```

Generate multiple documents with the details.

```bash
pgreport "SELECT 'communities/' || bfs_id || '.json' AS filename, \\
                 json_agg(c) AS document \\
          FROM communities) AS c"
```

## Generate CSV files

Create a single TSV file containing all flat data. You cannot represent
structured data in TSV files. You can fallback to create hierarchies
using different files.

`pgreport` will automatically detect that you want to create a TSV file and
will choose sensible defaults for you.

```bash
pgreport "SELECT 'communities.tsv' AS filename, \\
                 bfs_id, name \\
          FROM communities"
```

## Generate XML files

But XML is dead? Many applications still prefer XML as a data format and if you don't
have to support a specific schema or want to get input for XSLT `pgreport` can generate
the necessary files for you. You can either rely on default XML output
or build your own XML document with [XML functions in PostgreSQL](https://wiki.postgresql.org/wiki/XML_Support).

```bash
pgreport "SELECT 'communities.tsv' AS filename, \\
                 bfs_id, name \\
          FROM communities"
```



## Install

You can download a single binary for Linux, OSX or Windows.

**OSX**

```bash
wget -O pgreport https://github.com/lukasmartinelli/pgfutter/releases/download/v0.3.2/pgfutter_darwin_amd64
chmod +x pgreport

./pgreport --help
```

**Linux**

```bash
wget -O pgreport https://github.com/lukasmartinelli/pgfutter/releases/download/v0.3.2/pgfutter_linux_amd64
chmod +x pgreport

./pgreport --help
```

**Install from source**

```bash
go get github.com/lukasmartinelli/pgreport
```

If you are using Windows or 32-bit architectures you need to [download the appropriate binary
yourself](https://github.com/lukasmartinelli/pgreport/releases/latest).

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

## Personal Motivation

I use PostgreSQL in most ETL workflows to consolidate, aggregate and cleanup data.
After doing that I want to get the data out again which previously relied on
a lot of redundant Python code code projects all of which has now been replaces
with `pgreport`.

## Advanced Use Cases

### Load SQL from File

If you have a long SQL statement to select your data you can read
the query from a file. Instead of passing a query to `pgreport` you 
pass a filename ending with `.sql`.

```bash
# Store query in file
echo 'SELECT * FROM communities' > myquery.sql
# Execute query from file
pgreport jsonlines myquery.sql
```

## Alternatives

- [ ] Research alternatives

## Cross-compiling

We use [gox](https://github.com/mitchellh/gox) to create distributable
binaries for Windows, OSX and Linux.

```bash
docker run --rm -v "$(pwd)":/usr/src/pgreport -w /usr/src/pgreport tcnksm/gox:1.4.2-light
```
