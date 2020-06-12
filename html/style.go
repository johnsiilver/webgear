package html

import (
	"fmt"
	"html/template"
	"strings"
)

var styleTmpl = template.Must(template.New("style").Parse(strings.TrimSpace(`
<style {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
{{.Self.TagValue}}
</style>
`)))

// Style defines an HTML style tag.
type Style struct {
	GlobalAttrs

	// TagValue provides the value inside a reference.
	TagValue template.CSS

	Events *Events
}

func (s *Style) validate() error {
	if s.TagValue == "" {
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

func (s *Style) Execute(pipe Pipeline) string {
	pipe.Self = s

	if err := styleTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
