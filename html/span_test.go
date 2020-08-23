package html

import (
	"context"
	"strings"
	"testing"
)

func TestSpan(t *testing.T) {
	tests := []struct {
		desc string
		span *Span
		want string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			span: &Span{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Elements: []Element{TextElement("text")},
				Events:   (&Events{}).AddScript(OnError, "handleError"),
			},

			want: strings.TrimSpace(`
<span accesskey="key" onerror="handleError">
	text
</span>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.span.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestSpan(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
