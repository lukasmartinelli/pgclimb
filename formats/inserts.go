// Implements INSERTS output for exported rows
// e.g.
// INSERT INTO <table> (name, last_name, something);
package formats

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	insertQuery = "INSERT INTO %s (%s) "
)

// very simple quoting for sql values
func quote(val string) string {
	buf := bytes.NewBufferString("'")
	buf.WriteString(strings.Replace(val, "'", "\\'", -1))
	buf.WriteString("'")
	return buf.String()
}

type InsertsFormat struct {
	DataFormat
	tableName   string
	columns     []string
	multiInsert bool
	writer      *bufio.Writer
}

func NewInsertsFormat(w io.Writer, fileName string, tableName string) (*InsertsFormat, error) {
	return &InsertsFormat{
		writer:      bufio.NewWriter(w),
		tableName:   tableName,
		columns:     make([]string, 0),
		multiInsert: false,
	}, nil
}

func (f *InsertsFormat) WriteHeader(columns []string) error {
	f.columns = columns
	return nil
}

func (f *InsertsFormat) Flush() error { return nil }

func (f *InsertsFormat) WriteRow(values map[string]interface{}) error {
	columnsVal := strings.Join(f.columns, ",")
	queryStr := fmt.Sprintf(insertQuery, f.tableName, columnsVal)
	buf := bytes.NewBufferString(queryStr)
	record := []string{}
	for _, col := range f.columns {
		switch value := (values[col]).(type) {
		case string:
			record = append(record, quote(value))
		case []byte:
			record = append(record, quote(string(value)))
		case int64:
			record = append(record, fmt.Sprintf("%d", value))
		case float64:
			record = append(record, strconv.FormatFloat(value, 'f', -1, 64))
		case time.Time:
			record = append(record, quote(value.Format(time.RFC3339)))
		case bool:
			if value == true {
				record = append(record, "true")
			} else {
				record = append(record, "false")
			}
		case nil:
			record = append(record, "null")
		}
	}
	buf.WriteString("VALUES (")
	buf.WriteString(strings.Join(record, ","))
	buf.WriteString(");\n")

	_, err := f.writer.Write(buf.Bytes())
	if err != nil {
		return err
	}
	err = f.writer.Flush()
	// defer f.writer.Close()
	return err
}
