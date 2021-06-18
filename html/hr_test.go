package html

import (
	"context"
	"strings"
	"testing"
)

func TestHR(t *testing.T) {
	tests := []struct {
		desc string
		hr   *HR
		want string
	}{
		{
			desc: "Empty attributes",
			hr:   &HR{},
			want: "<hr  >",
		},
		{
			desc: "1 global + 1 event",
			hr: &HR{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).AddScript(OnError, "handleError"),
			},

			want: strings.TrimSpace(`<hr accesskey="key" onerror="handleError">`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.hr.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestHR(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
