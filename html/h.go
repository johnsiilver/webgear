package html

import (
	"fmt"
	"html/template"
	"strings"
)

var hTmpl = template.Must(template.New("h").Parse(strings.TrimSpace(`
<h{{.Self.Level}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</h{{.Self.Level}}>
`)))

// Level represents the level of the tag, 1-6.
type Level uint8

// H represents a H1-6 tag.
type H struct {
	GlobalAttrs

	// Level is the tag's level, 1-6.
	Level Level

	// Elements are elements contained within the H tag.
	Elements []Element

	Events *Events
}

func (h *H) validate() error {
	if uint8(h.Level) < 1 || uint8(h.Level) > 6 {
		return fmt.Errorf("H tag has level %d, must be 1 to 6", h.Level)
	}
	return nil
}

func (h *H) Execute(pipe Pipeline) string {
	pipe.Self = h

	if err := hTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
