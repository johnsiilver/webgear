package html

import (
	"strings"
	"testing"
)

func TestBody(t *testing.T) {
	tests := []struct {
		desc string
		body *Body
		want string
	}{
		{
			desc: "Empty attributes",
			body: &Body{},
			want: "<body  >\n</body>",
		},
		{
			desc: "All attributes + 1 global + 1 event + 1 element",
			body: &Body{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).OnError("handleError"),
				Elements: []Element{
					&Div{
						Elements: []Element{
							&A{Href: "/subpage", TagValue: TextElement("hello")},
						},
					},
				},
			},

			want: strings.TrimSpace(`
<body accesskey="key" onerror="handleError">
	<div  >
	<a href="/subpage"  >hello</a>
</div>
</body>
`),
		},
	}

	for _, test := range tests {
		if err := test.body.compile(); err != nil {
			panic(err)
		}
		got := test.body.Execute(struct{}{})
		if test.want != string(got) {
			t.Errorf("TestBody(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
