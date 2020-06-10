package html

import (
	"strings"
	"testing"
)

func TestIFrame(t *testing.T) {
	tests := []struct {
		desc   string
		iframe *IFrame
		want   string
	}{
		{
			desc: "Most attributes using Src + two sandboxing params + 1 global + 1 event",
			iframe: &IFrame{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events:              (&Events{}).OnError("handleError"),
				Name:                "name",
				Src:                 URLParse("https://vimeo.com"),
				Allow:               "autoplay; fullscreen",
				AllowFullscreen:     true,
				AllowPaymentRequest: true,
				Height:              110,
				Width:               110,
				ReferrerPolicy:      OriginWhenCrossOrigin,
				Sandboxing:          Sandboxing{AllowFormsSB, AllowPopupsSB},
			},

			want: strings.TrimSpace(`
<iframe name="name" src="https://vimeo.com" allow="autoplay; fullscreen" allowfullscreen="true" ` +
				`allowpaymentrequest="true" height="110" width="110" referrerpolicy="origin-when-cross-origin" ` +
				`sandbox="allow-forms allow-popups" accesskey="key" onerror="handleError"></iframe>
`),
		},
	}

	for _, test := range tests {
		if err := test.iframe.Init(); err != nil {
			panic(err)
		}
		got := test.iframe.Execute(Pipeline{})
		if test.want != string(got) {
			t.Errorf("TestIFrame(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
