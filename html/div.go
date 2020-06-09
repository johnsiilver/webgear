package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var divTmpl = strings.TrimSpace(`
<div {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</div>
`)

// Div represents a division tag.
type Div struct {
	GlobalAttrs
	// Elements are elements contained within the Div.
	Elements []Element

	Events *Events

	tmpl *template.Template

	pool sync.Pool
}

func (d *Div) isElement() {}

func (d *Div) compile() error {
	var err error
	d.tmpl, err = template.New("div").Parse(divTmpl)
	if err != nil {
		return fmt.Errorf("Div object had error: %s", err)
	}

	for _, element := range d.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

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

	if err := d.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
