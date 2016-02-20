package formats

import (
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
)

type XlsxFormat struct {
	file     *xlsx.File
	sheet    *xlsx.Sheet
	fileName string
	columns  []string
}

func NewXlsxFormat(fileName string, sheetName string) (*XlsxFormat, error) {
	file := xlsx.NewFile()
	if _, err := os.Stat(fileName); fileName != "" && err == nil {
		file, err = xlsx.OpenFile(fileName)
		if err != nil {
			fmt.Println("Errord file")
			return nil, err
		}
	}

	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		// Sheet already exists - empty it first
		sheet = file.Sheet[sheetName]
		sheet.Rows = make([]*xlsx.Row, 0)
	}

	return &XlsxFormat{
		file:     file,
		fileName: fileName,
		sheet:    sheet,
		columns:  make([]string, 0),
	}, nil
}

func (f *XlsxFormat) Flush() error {
	if f.fileName == "" {
		return f.file.Write(os.Stdout)
	} else {
		return f.file.Save(f.fileName)
	}
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
