package html

import (
	"html/template"
	"strings"
)

var brTmpl = template.Must(template.New("h").Parse(strings.TrimSpace(`
<br {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
`)))

// BR represents line break tag.
type BR struct {
	GlobalAttrs

	Events *Events
}

func (b *BR) Execute(pipe Pipeline) string {
	pipe.Self = b

	if err := brTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
