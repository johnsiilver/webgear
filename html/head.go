package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var headTmpl = strings.TrimSpace(`
<head {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .Data}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</head>
`)

// Head represents an HTML head tag.
type Head struct {
	GlobalAttrs

	// Elements are elements contained within the Head.
	Elements []Element

	Events *Events

	pool sync.Pool

	tmpl *template.Template
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

func (h *Head) compile() error {
	var err error
	h.tmpl, err = template.New("head").Parse(headTmpl)
	if err != nil {
		return fmt.Errorf("Head object had error: %s", err)
	}

	for _, element := range h.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

	h.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (h *Head) Execute(data interface{}) template.HTML {
	buff := h.pool.Get().(*strings.Builder)
	defer h.pool.Put(buff)
	buff.Reset()

	if err := h.tmpl.Execute(buff, pipeline{Self: h, Data: data}); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}