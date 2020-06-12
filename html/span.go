package html

import (
	"html/template"
	"strings"
)

var spanTmpl = template.Must(template.New("span").Parse(strings.TrimSpace(`
<span {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</span>
`)))

// Span tag is an inline container used to mark up a part of a text, or a part of a document.
type Span struct {
	GlobalAttrs

	// Element is any containing element.
	Elements []Element

	Events *Events
}

func (s *Span) Execute(pipe Pipeline) string {
	pipe.Self = s

	if err := spanTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
