package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var bodyTmpl = template.Must(template.New("body").Parse(strings.TrimSpace(`
{{- if not .Self.Component}}<body {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>{{- end}}
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
{{if not .Self.Component -}}</body>{{- end}}
`)))

// Body represents the HTML body.
type Body struct {
	GlobalAttrs

	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events

	// Componenet is used to indicate that this is a snippet of code, not a full document.
	// As such, <body> will suppressed.
	Component bool

	pool sync.Pool
}

func (b *Body) Init() error {
	if err := compileElements(b.Elements); err != nil {
		return err
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

	if err := bodyTmpl.Execute(buff, pipe); err != nil {
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
