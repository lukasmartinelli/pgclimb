package main

import (
	"encoding/json"
	"os"
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

	supportsFilename := func() bool {
		for _, colName := range rc.colNames {
			if colName == "filename" {
				return true
			}
		}
		return false
	}()

	for rows.Next() {
		if err := rc.Update(rows); err != nil {
			return err
		}

		encoder := json.NewEncoder(os.Stdout)
		values := rc.Get()

		if supportsFilename {
			delete(values, "filename")
			encoder.Encode(values)
		} else {
			encoder.Encode(values)
		}
	}

	err = rows.Err()
	return err
}
