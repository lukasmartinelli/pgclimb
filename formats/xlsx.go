package formats

import (
	"io"

	"github.com/tealeg/xlsx"
)

type XlsxFormat struct {
	file    *xlsx.File
	sheet   *xlsx.Sheet
	writer  io.Writer
	columns []string
}

func NewXlsxFormat(w io.Writer, sheetName string) *XlsxFormat {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet(sheetName)

	return &XlsxFormat{
		file:    file,
		sheet:   sheet,
		writer:  w,
		columns: make([]string, 0),
	}
}

func (f *XlsxFormat) Flush() error {
	return f.file.Write(f.writer)
}

func (f *XlsxFormat) WriteHeader(columns []string) error {
	f.columns = columns
	row := f.sheet.AddRow()
	for _, col := range columns {
		cell := row.AddCell()
		cell.SetString(col)
	}
	return nil
}

func (f *XlsxFormat) WriteRow(values map[string]interface{}) error {
	row := f.sheet.AddRow()

	for _, col := range f.columns {
		cell := row.AddCell()
		switch value := (values[col]).(type) {
		case []byte:
			cell.SetString(string(value))
		case int64:
			cell.SetInt64(value)
		}
	}
	return nil
}
