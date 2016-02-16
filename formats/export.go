package formats

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lukasmartinelli/pgclimb/pg"
)

// Supports encoding a row in different formats
type RowEncoder interface {
	Encode(map[string]interface{}) error
}

func Export(query string, connStr string, encoder RowEncoder) error {
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

			if err = encoder.Encode(values); err != nil {
				return err
			}
			file.Sync()
			log.Printf("%s\n", filename)
		} else {
			if err = encoder.Encode(values); err != nil {
				return err
			}
		}
	}

	err = rows.Err()
	return err
}

func printValue(pval *interface{}) {
	//fmt.Println(reflect.ValueOf(*pval).Kind())
	switch v := (*pval).(type) {
	case nil:
		fmt.Print("NULL")
	case int64:
		fmt.Print(int64(v))
	case bool:
		if v {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	case []byte:
		fmt.Print(string(v))
	case time.Time:
		fmt.Print(v.Format("2006-01-02 15:04:05.999"))
	default:
		fmt.Print(v)
	}
}
