package html

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestImg(t *testing.T) {
	u, _ := url.Parse("/path")

	tests := []struct {
		desc string
		img  *Img
		want string
	}{
		{
			desc: "All attributes using size in px + 1 global + 1 event ",
			img: &Img{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events:         (&Events{}).AddScript(OnError, "handleError"),
				Src:            u,
				SrcSet:         u,
				Alt:            "alt",
				UseMap:         "#map",
				CrossOrigin:    UseCredentialsCO,
				HeightPx:       100,
				WidthPx:        100,
				IsMap:          true,
				LongDesc:       u,
				ReferrerPolicy: OriginWhenCrossOrigin,
				Sizes:          "sizes",
			},

			want: strings.TrimSpace(`
<img src="/path" srcset="/path" alt="alt" usemap="#map" crossorigin="use-credentials" ` +
				`height="100px" width="100px" ismap longdesc="/path" referrerpolicy="origin-when-cross-origin" sizes="sizes" ` +
				`accesskey="key" onerror="handleError"/>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.img.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestImg(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
