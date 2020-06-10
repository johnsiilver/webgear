package html

import (
	"html/template"
	"strings"
	"sync"
)

var divTmpl = template.Must(template.New("div").Parse(strings.TrimSpace(`
<div {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</div>
`)))

// Div represents a division tag.
type Div struct {
	GlobalAttrs
	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events

	pool sync.Pool
}

func (d *Div) Init() error {
	d.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (d *Div) Execute(pipe Pipeline) template.HTML {
	buff := d.pool.Get().(*strings.Builder)
	defer d.pool.Put(buff)
	buff.Reset()

	pipe.Self = d

	if err := divTmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
