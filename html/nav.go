package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var navTmpl = strings.TrimSpace(`
<nav {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</nav>
`)

// Nav defines an HTML nav tag.
type Nav struct {
	GlobalAttrs
	Events *Events

	Elements []Element

	tmpl *template.Template
	pool sync.Pool
}

func (n *Nav) isElement() {}

func (n *Nav) compile() error {
	var err error
	n.tmpl, err = template.New("nav").Parse(navTmpl)
	if err != nil {
		return fmt.Errorf("Ul object had error: %s", err)
	}

	for _, element := range n.Elements {
		if err := element.compile(); err != nil {
			return err
		}
	}

	n.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (n *Nav) Execute(pipe Pipeline) template.HTML {
	buff := n.pool.Get().(*strings.Builder)
	defer n.pool.Put(buff)
	buff.Reset()

	pipe.Self = n

	if err := n.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
