package html

import (
	"html/template"
	"strings"
)

var divTmpl = template.Must(template.New("div").Parse(strings.TrimSpace(`
<div {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</div>
`)))

// Div represents a division tag.
type Div struct {
	GlobalAttrs
	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events
}

func (d *Div) Execute(pipe Pipeline) string {
	pipe.Self = d

	if err := divTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

func (d *Div) isFormElement() {}
