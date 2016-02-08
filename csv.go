package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvEncoder struct {
	writer *csv.Writer
}

func NewCsvEncoder() *CsvEncoder {
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = '\t'
	return &CsvEncoder{writer}
}

func (e *CsvEncoder) Encode(values map[string]interface{}) error {
	record := []string{}
	for _, value := range values {
		var column string
		switch value := (value).(type) {
		case []byte:
			column = string(value)
		case int64:
			column = fmt.Sprintf("%d", value)
		}
		record = append(record, column)
	}
	err := e.writer.Write(record)
	if err != nil {
		return err
	}

	e.writer.Flush()
	err = e.writer.Error()
	return err
}
