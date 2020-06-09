package html

import (
	"strings"
	"testing"
)

func TestUl(t *testing.T) {
	tests := []struct {
		desc string
		ul   *Ul
		want string
	}{
		{
			desc: "Empty attributes",
			ul:   &Ul{},
			want: "<ul  >\n</ul>",
		},
		{
			desc: "All attributes + 1 global + 1 event + 1 element",
			ul: &Ul{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).OnError("handleError"),
				Elements: []Element{
					&A{Href: "/subpage", TagValue: TextElement("hello")},
				},
			},

			want: strings.TrimSpace(`
<ul accesskey="key" onerror="handleError">
	<a href="/subpage"  >hello</a>
</ul>
`),
		},
	}

	for _, test := range tests {
		if err := test.ul.compile(); err != nil {
			panic(err)
		}
		got := test.ul.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestNav(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
