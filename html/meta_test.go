package html

import (
	"context"
	"strings"
	"testing"
)

func TestMeta(t *testing.T) {
	tests := []struct {
		desc string
		meta *Meta
		want string
	}{

		{
			desc: "Charset + 1 global",
			meta: &Meta{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Charset: "UTF-8",
			},
			want: `<meta charset="UTF-8" accesskey="key">`,
		},
		{
			desc: "Charset + 1 global",
			meta: &Meta{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				MetaName: ApplicationNameMN,
				Content:  "content",
			},
			want: `<meta metaname="application-name" content="content" accesskey="key">`,
		},
		{
			desc: "Charset + 1 global",
			meta: &Meta{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				HTTPEquiv: ContentTypeHE,
				Content:   "content",
			},
			want: `<meta http-equiv="content-type" content="content" accesskey="key">`,
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.meta.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestMeta(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
