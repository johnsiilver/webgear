package html

import (
	"context"
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
				Events: (&Events{}).AddScript(OnError, "handleError"),
				Elements: []Element{
					&A{Href: "/subpage", Elements: []Element{TextElement("hello")}},
				},
			},

			want: strings.TrimSpace(`
<ul accesskey="key" onerror="handleError">
	<a href="/subpage"  >
	hello
</a>
</ul>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.ul.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestNav(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
