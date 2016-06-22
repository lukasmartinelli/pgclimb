package formats

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type CsvFormat struct {
	writer    *csv.Writer
	columns   []string
	headerRow bool
}

func NewCsvFormat(w io.Writer, delimiter rune, headerRow bool) *CsvFormat {
	writer := csv.NewWriter(w)
	writer.Comma = delimiter
	return &CsvFormat{
		writer:    writer,
		columns:   make([]string, 0),
		headerRow: headerRow,
	}
}

func (f *CsvFormat) WriteHeader(columns []string) error {
	f.columns = columns
	if f.headerRow {
		return f.writer.Write(columns)
	} else {
		return nil
	}
}

func (f *CsvFormat) Flush() error { return nil }

func (f *CsvFormat) WriteRow(values map[string]interface{}) error {
	record := []string{}
	for _, col := range f.columns {
		switch value := (values[col]).(type) {
		case string:
			record = append(record, value)
		case []byte:
			record = append(record, string(value))
		case int64:
			record = append(record, fmt.Sprintf("%d", value))
		case float64:
			record = append(record, strconv.FormatFloat(value, 'f', -1, 64))
		case time.Time:
			record = append(record, value.Format(time.RFC3339))
		case bool:
			if value == true {
				record = append(record, "true")
			} else {
				record = append(record, "false")
			}
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
