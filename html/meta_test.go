package html

import (
	"testing"
)

func TestMeta(t *testing.T) {
	tests := []struct {
		desc string
		meta    *Meta
		want string
	}{

		{
			desc: "Charset + 1 global",
			meta: &Meta{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Charset:	"UTF-8",
			},
			want: `<meta charset="UTF-8" accesskey="key">`,
		},
		{
			desc: "Charset + 1 global",
			meta: &Meta{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				MetaName:	ApplicationNameMN,
				Content: "content",
			},
			want: `<meta metaname="application-name" content="content" accesskey="key">`,
		},
		{
			desc: "Charset + 1 global",
			meta: &Meta{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				HTTPEquiv:	ContentTypeHE,
				Content: "content",
			},
			want: `<meta http-equiv="content-type" content="content" accesskey="key">`,
		},
	}

	for _, test := range tests {
		if err := test.meta.compile(); err != nil {
			panic(err)
		}
		got := test.meta.Execute(struct{}{})
		if test.want != string(got) {
			t.Errorf("TestMeta(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
