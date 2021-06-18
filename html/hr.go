package html

import (
	"html/template"
	"strings"
)

var hrTmpl = template.Must(template.New("hr").Parse(strings.TrimSpace(`<hr {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>`)))

// HR represents a horizontal rule tag.
type HR struct {
	GlobalAttrs
	// Events are events attached to the element.
	Events *Events
}

func (h *HR) Execute(pipe Pipeline) string {
	pipe.Self = h

	if err := hrTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
