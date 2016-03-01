package formats

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"time"
)

type XMLFormat struct {
	encoder *xml.Encoder
}

func NewXMLFormat(w io.Writer) *XMLFormat {
	e := xml.NewEncoder(w)
	e.Indent("  ", "    ")
	return &XMLFormat{e}
}

// Writing header for XML is a NOP
func (e *XMLFormat) WriteHeader(columns []string) error {
	return nil
}

func (e *XMLFormat) Flush() error { return nil }

func (e *XMLFormat) WriteRow(values map[string]interface{}) error {
	row := xml.StartElement{Name: xml.Name{"", "row"}}
	tokens := []xml.Token{row}
	for key, value := range values {
		var charData xml.CharData

		t := xml.StartElement{Name: xml.Name{"", key}}

		switch value := (value).(type) {
		case []byte:
			charData = xml.CharData(string(value))
		case int64:
			charData = xml.CharData(fmt.Sprintf("%d", value))
		case float64:
			charData = xml.CharData(strconv.FormatFloat(value, 'f', -1, 64))
		case time.Time:
			charData = xml.CharData(value.Format(time.RFC3339))
		case bool:
			if value == true {
				charData = xml.CharData("true")
			} else {
				charData = xml.CharData("false")
			}
		}
		tokens = append(tokens, t, charData, t.End())
	}
	tokens = append(tokens, row.End())

	for _, t := range tokens {
		err := e.encoder.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	err := e.encoder.Flush()
	if err != nil {
		return err
	}

	return nil
}
