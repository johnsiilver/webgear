package html

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
)

var styleTmpl = strings.TrimSpace(`
<style {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{.Self.TagValue}}
</style>
`)

// Style defines an HTML style tag.
type Style struct {
	GlobalAttrs

	// TagValue provides the value inside a reference.
	TagValue TextElement

	Events *Events

	tmpl *template.Template

	pool sync.Pool
}

func (s *Style) validate() error {
	if s.TagValue.isZero() {
		return fmt.Errorf("Style element cannot have a nil TagValue")
	}
	if strings.TrimSpace(string(s.TagValue)) == "" {
		return fmt.Errorf("Style isn't empty, but it only contains space characters, which is also invalid. Nice try")
	}
	return nil
}

func (s *Style) Attr() template.HTMLAttr {
	output := structToString(s)
	return template.HTMLAttr(output)
}

func (s *Style) isElement() {}

func (s *Style) compile() error {
	var err error
	s.tmpl, err = template.New("s").Parse(styleTmpl)
	if err != nil {
		return fmt.Errorf("Style object had error: %s", err)
	}

	s.pool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
	return nil
}

func (s *Style) Execute(pipe Pipeline) template.HTML {
	buff := s.pool.Get().(*strings.Builder)
	defer s.pool.Put(buff)
	buff.Reset()

	pipe.Self = s

	if err := s.tmpl.Execute(buff, pipe); err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}
