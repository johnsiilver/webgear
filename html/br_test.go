package html

import (
	"context"
	"strings"
	"testing"
)

func TestBR(t *testing.T) {
	tests := []struct {
		desc string
		br   *BR
		want string
	}{
		{
			desc: "Empty attributes",
			br:   &BR{},
			want: "<br  >",
		},
		{
			desc: "1 global + 1 event",
			br: &BR{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).AddScript(OnError, "handleError"),
			},

			want: strings.TrimSpace(`
<br accesskey="key" onerror="handleError">
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.br.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestBR(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
