package html

import (
	"html/template"
	"strings"
)

var navTmpl = template.Must(template.New("nav").Parse(strings.TrimSpace(`
<nav {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</nav>
`)))

// Nav defines an HTML nav tag.
type Nav struct {
	GlobalAttrs
	Events *Events

	Elements []Element
}

func (n *Nav) Execute(pipe Pipeline) string {
	pipe.Self = n

	if err := navTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
