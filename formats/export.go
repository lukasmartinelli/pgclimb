package formats

import "github.com/lukasmartinelli/pgclimb/pg"

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

	if err = format.WriteHeader(columnNames); err != nil {
		return err
	}

	for rows.Next() {
		values := make(map[string]interface{})
		if err = rows.MapScan(values); err != nil {
			return err
		}

		if err = format.WriteRow(values); err != nil {
			return err
		}
	}

	if err = format.Flush(); err != nil {
		return err
	}

	return rows.Err()
}
