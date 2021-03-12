package html

import (
	"context"
	"strings"
	"testing"
)

func TestLabel(t *testing.T) {
	tests := []struct {
		desc  string
		label *Label
		want  string
	}{
		{
			desc:  "Empty attributes",
			label: &Label{},
			want:  "<label  >\n</label>",
		},

		{
			desc: "All attributes + 1 global + 1 event + 1 element",
			label: &Label{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).AddScript(OnError, "handleError"),
				Elements: []Element{
					&A{Href: "/subpage", Elements: []Element{TextElement("hello")}},
				},
			},

			want: strings.TrimSpace(`
<label accesskey="key" onerror="handleError">
	<a href="/subpage"  >
	hello
</a>
</label>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.label.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestLabel(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
