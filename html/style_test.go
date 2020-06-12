package html

import (
	"context"
	"html/template"
	"strings"
	"testing"
)

func TestStyle(t *testing.T) {
	tests := []struct {
		desc  string
		style *Style
		want  string
	}{
		{
			desc:  "Empty attributes",
			style: &Style{},
			want: strings.TrimSpace(`
<style  >

</style>
`),
		},
		{
			desc: "All attributes + 1 global + 1 event",
			style: &Style{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				TagValue: template.CSS("text"),
				Events:   (&Events{}).OnError("handleError"),
			},
			want: strings.TrimSpace(
				`<style accesskey="key" onerror="handleError">
text
</style>`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.style.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestStyle(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
