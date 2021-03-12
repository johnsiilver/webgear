package html

import (
	"html/template"
	"strings"
)

var iTmpl = template.Must(template.New("i").Parse(strings.TrimSpace(`
<i {{.Self.Attr }} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- .Self.Element.Execute $data}}
</i>
`)))

// I defines a part of text in an alternate voice or mood.
type I struct {
	GlobalAttrs

	Element TextElement

	Events *Events
}

func (i *I) Attr() template.HTMLAttr {
	output := structToString(i)
	return template.HTMLAttr(output)
}

func (i *I) Execute(pipe Pipeline) string {
	pipe.Self = i

	if err := iTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
