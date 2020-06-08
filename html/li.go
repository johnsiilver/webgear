package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var liTmpl = strings.TrimSpace(`
<li {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .Data}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</li>
`)

// Li defines an HTML li tag.
type Li struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	tmpl *template.Template
	pool sync.Pool
}

func (l *Li) isElement() {}

func (l *Li) compile() error {
	var err error
	l.tmpl, err = template.New("li").Parse(liTmpl)
	if err != nil {
		return fmt.Errorf("Ul object had error: %s", err)
	}

	for _, element := range l.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

	l.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (l *Li) Execute(data interface{}) template.HTML {
	buff := l.pool.Get().(*strings.Builder)
	defer l.pool.Put(buff)
	buff.Reset()

	if err := l.tmpl.Execute(buff, pipeline{Self: l, Data: data}); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
