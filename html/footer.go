package html

import (
	"html/template"
)

type Footer struct {
}

func (f Footer) compile(tmpl *template.Template) error {
	tmpl.New("footer").Parse(``)
	return nil
}
