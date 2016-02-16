package formats

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CsvFormat struct {
	writer *csv.Writer
}

func NewCsvFormat(delimiter rune) *CsvFormat {
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = delimiter
	return &CsvFormat{writer}
}

func (e *CsvFormat) WriteHeader(columns []string) error {
	return e.writer.Write(columns)
}

func (e *CsvFormat) Flush() error { return nil }

func (e *CsvFormat) WriteRow(values map[string]interface{}) error {
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
