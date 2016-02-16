package formats

import (
	"log"
	"os"

	"github.com/lukasmartinelli/pgclimb/pg"
)

// Supports storing data in different formats
type DataFormat interface {
	WriteHeader(columns []string) error
	WriteRow(map[string]interface{}) error
	Flush() error
}

func Export(query string, connStr string, format DataFormat) error {
	db, err := pg.Connect(connStr)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Queryx(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	supportsFilename := func() bool {
		for _, colName := range columnNames {
			if colName == "filename" {
				return true
			}
		}
		return false
	}()

	format.WriteHeader(columnNames)

	for rows.Next() {
		values := make(map[string]interface{})
		if err = rows.MapScan(values); err != nil {
			return nil
		}

		if supportsFilename {
			filename := values["filename"].(string)
			delete(values, "filename")

			file, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer file.Close()

			if err = format.WriteRow(values); err != nil {
				return err
			}
			file.Sync()
			log.Printf("%s\n", filename)
		} else {
			if err = format.WriteRow(values); err != nil {
				return err
			}
		}
	}

	if err = format.Flush(); err != nil {
		return err
	}

	err = rows.Err()
	return err
}
