package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var ulTmpl = strings.TrimSpace(`
<ul {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</ul>
`)

// Ul defines an HTML ul tag.
type Ul struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	tmpl *template.Template
	pool sync.Pool
}

func (u *Ul) isElement() {}

func (u *Ul) compile() error {
	var err error
	u.tmpl, err = template.New("ul").Parse(ulTmpl)
	if err != nil {
		return fmt.Errorf("Ul object had error: %s", err)
	}

	for _, element := range u.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

	u.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (u *Ul) Execute(pipe Pipeline) template.HTML {
	buff := u.pool.Get().(*strings.Builder)
	defer u.pool.Put(buff)
	buff.Reset()

	pipe.Self = u

	if err := u.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
