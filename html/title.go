package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var titleTmpl = strings.TrimSpace(`
<title {{.Self.GlobalAttrs.Attr}}>{{.Self.TagValue}}</title>
`)

// A defines a hyperlink, which is used to link from one page to another.
type Title struct {
	GlobalAttrs

	// TagValue provides the value inside a reference.
	TagValue TextElement

	tmpl *template.Template

	pool sync.Pool
}

func (t *Title) isElement() {}

func (t *Title) validate() error {
	if t.TagValue.isZero() {
		return fmt.Errorf("Title element cannot have a nil TagValue")
	}
	if strings.TrimSpace(string(t.TagValue)) == "" {
		return fmt.Errorf("Title isn't empty, but it only contains space characters, which is also invalid. Nice try")
	}
	return nil
}

func (t *Title) compile() error {
	var err error
	t.tmpl, err = template.New("title").Parse(titleTmpl)
	if err != nil {
		return fmt.Errorf("A object had error: %s", err)
	}

	t.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}

	return nil
}

func (t *Title) Execute(pipe Pipeline) template.HTML {
	buff := t.pool.Get().(*strings.Builder)
	defer t.pool.Put(buff)
	buff.Reset()

	pipe.Self = t

	if err := t.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
