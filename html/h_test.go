package html

import (
	"context"
	"strings"
	"testing"
)

func TestH(t *testing.T) {
	tests := []struct {
		desc string
		h    *H
		want string
	}{
		{
			desc: "Empty attributes",
			h:    &H{Level: 1},
			want: "<h1  >\n</h1>",
		},
		{
			desc: "1 element + 1 global + 1 event",
			h: &H{
				Level: 2,
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events:   (&Events{}).OnError("handleError"),
				Elements: []Element{TextElement("hello")},
			},

			want: strings.TrimSpace(`
<h2 accesskey="key" onerror="handleError">
	hello
</h2>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.h.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestH(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
