package html

import (
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
					&A{Href: "/subpage", TagValue: TextElement("hello")},
				},
			},

			want: strings.TrimSpace(`
<div accesskey="key" onerror="handleError">
	<a href="/subpage"  >hello</a>
</div>
`),
		},
	}

	for _, test := range tests {
		if err := test.div.compile(); err != nil {
			panic(err)
		}
		got := test.div.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestDiv(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
