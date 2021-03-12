package html

import (
	"context"
	"html/template"
	"strings"
	"testing"
)

type fakeGear struct {
	GearType
	Element
	name string
}

func (f fakeGear) TagType() template.HTMLAttr {
	return template.HTMLAttr(f.name)
}

func (f fakeGear) Execute(pipe Pipeline) string {
	return EmptyString
}

func TestComponent(t *testing.T) {
	tests := []struct {
		desc      string
		component *Component
		want      string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			component: &Component{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Gear:     fakeGear{name: "myComponent"},
				TagValue: TextElement("value"),
				Events:   (&Events{}).AddScript(OnError, "handleError"),
			},

			want: strings.TrimSpace(`
<myComponent accesskey="key" id="myComponent" onerror="handleError">
	value
</myComponent>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.component.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestComponent(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
