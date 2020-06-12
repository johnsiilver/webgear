package html

import (
	"fmt"
	"html/template"
	"strings"
)

var headTmpl = template.Must(template.New("head").Parse(strings.TrimSpace(`
<head {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{$data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</head>
`)))

// Head represents an HTML head tag.
type Head struct {
	GlobalAttrs

	// Elements are elements contained within the Head.
	Elements []Element

	Events *Events
}

func (h *Head) validate() error {
	hasTitle := false
	baseCount := 0
	for _, e := range h.Elements {
		switch e.(type) {
		case *Title:
			hasTitle = true
			continue
		case *Base:
			baseCount++
			continue
		case *Style, *Meta, *Link, *Script:
			continue
		}
		return fmt.Errorf("Head element contained element type(%T) that was not Title, Style, Meta, Link, Script or Base", e)
	}

	switch false {
	case hasTitle:
		return fmt.Errorf("Head element has no Title element. HTML 5 spec says this is required")
	}

	if baseCount > 1 {
		return fmt.Errorf("Head element can have 0 or 1 Base elements")
	}

	return nil
}

func (h *Head) Init() error {
	if err := compileElements(h.Elements); err != nil {
		return err
	}

	return nil
}

func (h *Head) Execute(pipe Pipeline) string {
	pipe.Self = h

	if err := headTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
