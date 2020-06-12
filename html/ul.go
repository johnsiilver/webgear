package html

import (
	"html/template"
	"strings"
)

var ulTmpl = template.Must(template.New("ul").Parse(strings.TrimSpace(`
<ul {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</ul>
`)))

// Ul defines an HTML ul tag.
type Ul struct {
	GlobalAttrs
	Events *Events

	Elements []Element
}

func (u *Ul) Execute(pipe Pipeline) string {
	pipe.Self = u

	if err := ulTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
