package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type RowEncoder interface {
	Encode(*io.Writer, map[string]string) error
}

// Try to JSON decode the bytes
func tryUnmarshal(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	return err
}

func export(query string, connStr string) error {
	db, err := connect(connStr)
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

			encoder := json.NewEncoder(file)
			encoder.Encode(values)
			file.Sync()
			log.Printf("%s\n", filename)
		} else {
			encoder := json.NewEncoder(os.Stdout)
			encoder.Encode(convertBytesToString(values))
			/*for _, val := range values {
				printValue(&val)
			}
			fmt.Println()
			*/
		}
	}

	err = rows.Err()
	return err
}

func convertBytesToString(rows map[string]interface{}) map[string]interface{} {
	for k, v := range rows {
		switch v := (v).(type) {
		case []byte:
			rows[k] = string(v)
		default:
			rows[k] = v
		}
	}
	return rows
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
