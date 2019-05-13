package core

import (
	"html/template"
	"io"
)

type Report struct {
	Session  *Session
	Template string
}

func (r *Report) Render(dest io.Writer) error {
	funcMap := template.FuncMap{
		"json": func(json string) template.JS {
			return template.JS(json)
		},
	}

	tmpl, err := template.New("Aquatone Report").Funcs(funcMap).Parse(r.Template)
	if err != nil {
		return err
	}

	err = tmpl.Execute(dest, r.Session)
	if err != nil {
		return err
	}

	return nil
}

func NewReport(s *Session, templ string) *Report {
	return &Report{
		Session:  s,
		Template: templ,
	}
}
