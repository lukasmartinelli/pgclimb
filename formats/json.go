package formats

import (
	"encoding/json"
	"os"
)

type JSONFormat struct {
	encoder *json.Encoder
}

func NewJSONFormat() *JSONFormat {
	return &JSONFormat{json.NewEncoder(os.Stdout)}
}

// Writing header for JSON is a NOP
func (e *JSONFormat) WriteHeader(columns []string) error {
	return nil
}

func (e *JSONFormat) WriteRow(rows map[string]interface{}) error {
	rows = convertBytesToString(rows)
	err := e.encoder.Encode(rows)
	return err
}

func convertBytesToString(rows map[string]interface{}) map[string]interface{} {
	for k, v := range rows {
		switch v := (v).(type) {
		case []byte:
			rows[k] = string(v)
		default:
			rows[k] = v
		}
	}
	return rows
}
