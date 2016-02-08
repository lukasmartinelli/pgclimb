package main

import (
	"database/sql"
	"fmt"

	"github.com/codegangsta/cli"
)

//setup a database connection and create the import schema
func connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
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

// Selecting unknown amount of columns from query
// https://gist.github.com/SchumacherFM/69a167bec7dea644a20e
type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
}

func NewMapStringScan(columnNames []string) *mapStringScan {
	lenCN := len(columnNames)
	s := &mapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}

func (s *mapStringScan) Get() map[string]string {
	return s.row
}
