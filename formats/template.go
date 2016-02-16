package formats

import (
	"os"
	"text/template"
)

type TemplateFormat struct {
	rows     []map[string]interface{}
	template *template.Template
}

func NewTemplateFormat(rawTemplate string) *TemplateFormat {
	t := template.Must(template.New("climbtemplate").Parse(rawTemplate))
	return &TemplateFormat{make([]map[string]interface{}, 0), t}
}

func (e *TemplateFormat) Flush() error {
	err := e.template.Execute(os.Stdout, e.rows)
	return err
}

func (e *TemplateFormat) WriteHeader(columns []string) error {
	return nil
}

func (e *TemplateFormat) WriteRow(values map[string]interface{}) error {
	e.rows = append(e.rows, convertToJSON(values))
	return nil
}
