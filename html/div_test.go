package html

import (
	"context"
	"strings"
	"testing"
)

func TestDiv(t *testing.T) {
	tests := []struct {
		desc string
		div  *Div
		want string
	}{
		{
			desc: "Empty attributes",
			div:  &Div{},
			want: "<div  >\n</div>",
		},
		{
			desc: "All attributes + 1 global + 1 event + 1 element",
			div: &Div{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).OnError("handleError"),
				Elements: []Element{
					&A{Href: "/subpage", Elements: []Element{TextElement("hello")}},
				},
			},

			want: strings.TrimSpace(`
<div accesskey="key" onerror="handleError">
	<a href="/subpage"  >
	hello
</a>
</div>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.div.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestDiv(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
