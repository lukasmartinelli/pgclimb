package formats

import (
	"encoding/json"
	"os"
)

type JsonEncoder struct {
	encoder *json.Encoder
}

func NewJsonEncoder() *JsonEncoder {
	return &JsonEncoder{json.NewEncoder(os.Stdout)}
}

func (e *JsonEncoder) Encode(rows map[string]interface{}) error {
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
