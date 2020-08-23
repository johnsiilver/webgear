package html

import (
	"context"
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
				Events: (&Events{}).AddScript(OnError, "handleError"),
				Elements: []Element{
					&Div{
						Elements: []Element{
							&A{Href: "/subpage", Elements: []Element{TextElement("hello")}},
						},
					},
				},
			},

			want: strings.TrimSpace(`
<body accesskey="key" onerror="handleError">
	<div  >
	<a href="/subpage"  >
	hello
</a>
</div>
</body>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.body.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestBody(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
