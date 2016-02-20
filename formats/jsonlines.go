package formats

import (
	"encoding/json"
	"io"
)

type JSONLinesFormat struct {
	encoder *json.Encoder
}

func NewJSONLinesFormat(w io.Writer) *JSONLinesFormat {
	return &JSONLinesFormat{json.NewEncoder(w)}
}

// Writing header for JSON is a NOP
func (e *JSONLinesFormat) WriteHeader(columns []string) error {
	return nil
}

func (e *JSONLinesFormat) Flush() error { return nil }

func (e *JSONLinesFormat) WriteRow(rows map[string]interface{}) error {
	rows = convertToJSON(rows)
	err := e.encoder.Encode(rows)
	return err
}
