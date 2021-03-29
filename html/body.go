package html

import (
	"fmt"
	"html/template"
	"strings"
)

var bodyTmpl = template.Must(template.New("body").Parse(strings.TrimSpace(`
{{- if not .Self.Component}}<body {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{- end}}
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
{{if not .Self.Component -}}</body>{{- end}}
`)))

// Body represents the HTML body.
type Body struct {
	GlobalAttrs

	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events

	// Componenet is used to indicate that this is a snippet of code, not a full document.
	// As such, the <body></body> tags will suppressed but the content will be rendered.
	Component bool
}

func (b *Body) Init() error {
	if err := compileElements(b.Elements); err != nil {
		return err
	}

	return nil
}

func (b *Body) Execute(pipe Pipeline) string {
	pipe.Self = b

	if err := bodyTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (b *Body) validate() error {
	if b == nil {
		return fmt.Errorf("Body element is not defined")
	}
	for _, e := range b.Elements {
		v, ok := e.(validator)
		if ok {
			if err := v.validate(); err != nil {
				return err
			}
		}
	}
	return nil
}
