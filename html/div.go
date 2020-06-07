package html

import (
	"fmt"
	"html/template"
	"strings"
)

var divTmpl = strings.TrimSpace(`
<div {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .Data}}
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

	str string
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

	return nil
}

func (d *Div) Execute(data interface{}) template.HTML {
	if d.str != "" {
		return template.HTML(d.str)
	}

	buff := strings.Builder{}

	if err := d.tmpl.Execute(&buff, pipeline{Self: d, Data: data}); err != nil {
		panic(err)
	}

	d.str = buff.String()
	return template.HTML(d.str)
}