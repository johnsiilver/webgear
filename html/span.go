package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var spanTmpl = strings.TrimSpace(`
<span {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</span>
`)

// Span tag is an inline container used to mark up a part of a text, or a part of a document.
type Span struct {
	GlobalAttrs

	// Element is any containing element.
	Elements []Element

	Events *Events

	tmpl *template.Template

	pool sync.Pool
}

func (s *Span) isElement() {}

func (s *Span) compile() error {
	var err error
	s.tmpl, err = template.New("s").Parse(spanTmpl)
	if err != nil {
		return fmt.Errorf("Span object had error: %s", err)
	}

	s.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (s *Span) Execute(pipe Pipeline) template.HTML {
	buff := s.pool.Get().(*strings.Builder)
	defer s.pool.Put(buff)
	buff.Reset()

	if err := s.tmpl.Execute(buff, Pipeline{Self: s}); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
