package formats

import (
	"github.com/tealeg/xlsx"

	"fmt"
)

type XlsxEncoder struct {
	file  *xlsx.File
	sheet *xlsx.Sheet
}

func NewXlsxEncoder() *XlsxEncoder {

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("data")

	return &XlsxEncoder{file, sheet}
}

func (e *XlsxEncoder) Encode(values map[string]interface{}) error {
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
