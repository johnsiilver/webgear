package html

import (
	"context"
	"strings"
	"testing"
	"html/template"
)

type fakeComponent struct {
	Element
	name string
}

func (f fakeComponent) TagType() template.HTMLAttr {
	return template.HTMLAttr(f.name)
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
				Gear: fakeComponent{name: "myComponent"},
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
