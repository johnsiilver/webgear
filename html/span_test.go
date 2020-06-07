package html

import (
	"testing"
	"strings"
)

func TestSpan(t *testing.T) {
	tests := []struct {
		desc string
		span    *Span
		want string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			span: &Span{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Element: TextElement("text"),
				Events:         (&Events{}).OnError("handleError"),
			},

			want: strings.TrimSpace(`
<span accesskey="key" onerror="handleError">
	text
</span>
`),
		},
	}

	for _, test := range tests {
		if err := test.span.compile(); err != nil {
			panic(err)
		}
		got := test.span.Execute(struct{}{})
		if test.want != string(got) {
			t.Errorf("TestSpan(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
