package formats

import (
	"io"
	"text/template"
)

type TemplateFormat struct {
	rows     []map[string]interface{}
	template *template.Template
	writer   io.Writer
}

func NewTemplateFormat(w io.Writer, rawTemplate string) *TemplateFormat {
	t := template.Must(template.New("climbtemplate").Parse(rawTemplate))
	return &TemplateFormat{
		rows:     make([]map[string]interface{}, 0),
		template: t,
		writer:   w,
	}
}

func (e *TemplateFormat) Flush() error {
	err := e.template.Execute(e.writer, e.rows)
	return err
}

func (e *TemplateFormat) WriteHeader(columns []string) error {
	return nil
}

func (e *TemplateFormat) WriteRow(values map[string]interface{}) error {
	e.rows = append(e.rows, convertToJSON(values))
	return nil
}
