package html

import (
	"fmt"
	"html/template"
	"strings"
)

var bodyTmpl = strings.TrimSpace(`
{{- if not .Self.Component}}<body {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{- end}}
	{{- $data := .Data}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
{{if not .Self.Component -}}</body>{{- end}}
`)

// Body represents the HTML body.
type Body struct {
	GlobalAttrs

	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events

	// Componenet is used to indicate that this is a snippet of code, not a full document.
	// As such, <body> will suppressed.
	Component bool

	tmpl *template.Template

	str string
}

func (b *Body) compile() error {
	var err error
	b.tmpl, err = template.New("body").Parse(bodyTmpl)
	if err != nil {
		return fmt.Errorf("Body object had error: %s", err)
	}

	for _, element := range b.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Body) Execute(data interface{}) template.HTML {
	if b.str != "" {
		return template.HTML(b.str)
	}

	buff := strings.Builder{}

	if err := b.tmpl.Execute(&buff, pipeline{Self: b, Data: data}); err != nil {
		panic(err)
	}

	b.str = buff.String()
	return template.HTML(b.str)
}

func (b *Body) validate() error {
	if b == nil {
		return fmt.Errorf("Body element is not defined")
	}
	return nil
}
