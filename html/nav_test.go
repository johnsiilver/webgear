package html

import (
	"strings"
	"testing"
)

func TestNav(t *testing.T) {
	tests := []struct {
		desc string
		nav  *Nav
		want string
	}{
		{
			desc: "Empty attributes",
			nav:  &Nav{},
			want: "<nav  >\n</nav>",
		},
		{
			desc: "All attributes + 1 global + 1 event + 1 element",
			nav: &Nav{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).OnError("handleError"),
				Elements: []Element{
					&A{Href: "/subpage", Elements: []Element{TextElement("hello")}},
				},
			},

			want: strings.TrimSpace(`
<nav accesskey="key" onerror="handleError">
	<a href="/subpage"  >
	hello
</a>
</nav>
`),
		},
	}

	for _, test := range tests {
		if err := test.nav.Init(); err != nil {
			panic(err)
		}
		if err := compileElement(test.nav); err != nil {
			panic(err)
		}
		got := test.nav.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestNav(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
