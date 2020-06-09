package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var bodyTmpl = strings.TrimSpace(`
{{- if not .Self.Component}}<body {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{- end}}
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
{{if not .Self.Component -}}</body>{{- end}}
`)

// Body represents the HTML body.
type Body struct {
	GlobalAttrs

	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events

	// Componenet is used to indicate that this is a snippet of code, not a full document.
	// As such, <body> will suppressed.
	Component bool

	tmpl *template.Template

	pool sync.Pool
}

func (b *Body) compile() error {
	var err error
	b.tmpl, err = template.New("body").Parse(bodyTmpl)
	if err != nil {
		return fmt.Errorf("Body object had error: %s", err)
	}

	for _, element := range b.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

	b.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (b *Body) Execute(pipe Pipeline) template.HTML {
	buff := b.pool.Get().(*strings.Builder)
	defer b.pool.Put(buff)
	buff.Reset()

	pipe.Self = b

	if err := b.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}

func (b *Body) validate() error {
	if b == nil {
		return fmt.Errorf("Body element is not defined")
	}
	return nil
}
