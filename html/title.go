package html

import (
	"fmt"
	"html/template"
	"strings"
)

var titleTmpl = template.Must(template.New("title").Parse(strings.TrimSpace(`
<title {{.Self.GlobalAttrs.Attr}}>{{.Self.TagValue}}</title>
`)))

// A defines a hyperlink, which is used to link from one page to another.
type Title struct {
	GlobalAttrs

	// TagValue provides the value inside a reference.
	TagValue TextElement
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

func (t *Title) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := titleTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
