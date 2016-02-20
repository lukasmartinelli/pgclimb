package formats

import (
	"encoding/json"
	"os"
)

type JSONArrayFormat struct {
	//Rows are stored in Memory until they are serialized into one Document
	//this only works because JSON documents are not supposed to be big
	//which would make them complicated to parse as well
	rows    []map[string]interface{}
	encoder *json.Encoder
}

func NewJSONArrayFormat() *JSONArrayFormat {
	return &JSONArrayFormat{make([]map[string]interface{}, 0), json.NewEncoder(os.Stdout)}
}

// Writing header for JSON is a NOP
func (e *JSONArrayFormat) WriteHeader(columns []string) error {
	return nil
}

func (e *JSONArrayFormat) Flush() error {
	err := e.encoder.Encode(e.rows)
	return err
}

func (e *JSONArrayFormat) WriteRow(rows map[string]interface{}) error {
	e.rows = append(e.rows, convertToJSON(rows))
	return nil
}

func convertToJSON(rows map[string]interface{}) map[string]interface{} {
	for k, v := range rows {
		switch v := (v).(type) {
		case []byte:
			var jsonVal interface{}
			err := json.Unmarshal(v, &jsonVal)
			if err == nil {
				rows[k] = jsonVal
			} else {
				rows[k] = string(v)
			}
		default:
			rows[k] = v
		}
	}
	return rows
}

// Try to JSON decode the bytes
func tryUnmarshal(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
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
