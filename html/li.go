package html

import (
	"html/template"
	"strings"
)

var liTmpl = template.Must(template.New("li").Parse(strings.TrimSpace(`
<li {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</li>
`)))

// Li defines an HTML li tag.
type Li struct {
	GlobalAttrs
	Events *Events

	Elements []Element
}

func (l *Li) Execute(pipe Pipeline) string {
	pipe.Self = l

	if err := liTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
