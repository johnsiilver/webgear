package html

import (
	"context"
	"strings"
	"testing"
)

func TestLi(t *testing.T) {
	tests := []struct {
		desc string
		li   *Li
		want string
	}{
		{
			desc: "Empty attributes",
			li:   &Li{},
			want: "<li  >\n</li>",
		},
		{
			desc: "All attributes + 1 global + 1 event + 1 element",
			li: &Li{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).AddScript(OnError, "handleError"),
				Elements: []Element{
					&A{Href: URLParse("/subpage"), Elements: []Element{TextElement("hello")}},
				},
			},

			want: strings.TrimSpace(`
<li accesskey="key" onerror="handleError">
	<a href="/subpage"  >
	hello
</a>
</li>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.li.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestLi(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
