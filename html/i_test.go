package html

import (
	"context"
	"strings"
	"testing"
)

func TestI(t *testing.T) {
	tests := []struct {
		desc string
		i    *I
		want string
	}{
		{
			desc: "Empty attributes",
			i:    &I{},
			want: `<i   >
</i>`,
		},
		{
			desc: "All attributes + 1 global + 1 event",
			i: &I{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Element:       TextElement("text"),
				Events:         (&Events{}).AddScript(OnError, "handleError"),
			},
			want: `<i  accesskey="key" onerror="handleError">text
</i>`,
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.i.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestI(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
