package html

import (
	"fmt"
	"sync"
	"strings"
	"html/template"
)

var spanTmpl = strings.TrimSpace(`
<span {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{.Self.Element.Execute .Data}}
</span>
`)

// Span tag is an inline container used to mark up a part of a text, or a part of a document.
type Span struct {
	GlobalAttrs

	// Element is any containing element.
	Element Element

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

func (s *Span) Execute(data interface{}) template.HTML {
	buff := s.pool.Get().(*strings.Builder)
	defer s.pool.Put(buff)
	buff.Reset()

	if err := s.tmpl.Execute(buff, pipeline{Self: s, Data: data}); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}