# pgreport ![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)

<img align="right" alt="Climbing elephant" src="logo.png" />

Generate JSON and CSV reports from PostgreSQL to create readonly APIs
and publish open data sets.

## Example

Let's generate a `communities.json` files containing an overview of all
files and a file for each community containing the details.

Generate a single document.

```bash
pgreport json "SELECT 'communities.json' AS filename, json_agg(t) AS document \\
                FROM (SELECT bfs_id, name FROM communities) AS t"
```

Generate multiple documents.

```bash
pgreport json "SELECT 'communities/' || bfs_id || '.json' AS filename, \\
                      json_agg(c) AS document \\
                FROM communities) AS c"
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

## Advanced Use Cases

### Custom delimiter

Quite often you want to specify a custom delimiter (default: `,`).

```bash
pgfutter csv -d "\t" traffic_violations.csv
```

You have to use `"` as a quoting character and `\` as escape character.
You might omit the quoting character if it is not necessary.

## Alternatives

- [ ] Research alternatives


## Cross-compiling

We use [gox](https://github.com/mitchellh/gox) to create distributable
binaries for Windows, OSX and Linux.

```bash
docker run --rm -v "$(pwd)":/usr/src/pgreport -w /usr/src/pgreport tcnksm/gox:1.4.2-light
```
