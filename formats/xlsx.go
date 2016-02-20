package formats

import (
	"fmt"
	"io"

	"github.com/tealeg/xlsx"
)

type XlsxFormat struct {
	file  *xlsx.File
	sheet *xlsx.Sheet
}

func NewXlsxFormat(w io.Writer) *XlsxFormat {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("data")

	return &XlsxFormat{file, sheet}
}

func (e *XlsxFormat) Flush() error { return nil }

func (e *XlsxFormat) WriteHeader(columns []string) error {
	row := e.sheet.AddRow()
	for _, col := range columns {
		cell := row.AddCell()
		cell.Value = col
	}
	return nil
}

func (e *XlsxFormat) WriteRow(values map[string]interface{}) error {
	row := e.sheet.AddRow()

	for _, value := range values {
		cell := row.AddCell()
		switch value := (value).(type) {
		case []byte:
			cell.Value = string(value)
		case int64:
			cell.Value = fmt.Sprintf("%d", value)
		}
	}
	return nil
}
