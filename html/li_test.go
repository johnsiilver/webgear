package html

import (
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
				Events: (&Events{}).OnError("handleError"),
				Elements: []Element{
					&A{Href: "/subpage", Elements: []Element{TextElement("hello")}},
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
		if err := test.li.Init(); err != nil {
			panic(err)
		}
		if err := compileElement(test.li); err != nil {
			panic(err)
		}

		got := test.li.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestLi(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
