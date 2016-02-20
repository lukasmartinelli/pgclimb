package formats

import (
	"encoding/csv"
	"fmt"
	"io"
)

type CsvFormat struct {
	writer  *csv.Writer
	columns []string
}

func NewCsvFormat(w io.Writer, delimiter rune) *CsvFormat {
	writer := csv.NewWriter(w)
	writer.Comma = delimiter
	return &CsvFormat{writer, make([]string, 0)}
}

func (f *CsvFormat) WriteHeader(columns []string) error {
	f.columns = columns
	return f.writer.Write(columns)
}

func (f *CsvFormat) Flush() error { return nil }

func (f *CsvFormat) WriteRow(values map[string]interface{}) error {
	record := []string{}
	for _, col := range f.columns {
		switch value := (values[col]).(type) {
		case []byte:
			record = append(record, string(value))
		case int64:
			record = append(record, fmt.Sprintf("%d", value))
		}
	}
	err := f.writer.Write(record)
	if err != nil {
		return err
	}

	f.writer.Flush()
	err = f.writer.Error()
	return err
}
