package main

import (
	"encoding/json"
	"log"
)

// Try to JSON decode the bytes
func tryUnmarshal(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	return err
}

func exportJSON(query string, connStr string) error {
	db, err := connect(connStr)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	rc := NewMapStringScan(columnNames)
	for rows.Next() {
		if err := rc.Update(rows); err != nil {
			return err
		}
		cv := rc.Get()
		log.Printf("%#v\n\n", cv)
	}

	err = rows.Err()
	return err
}
