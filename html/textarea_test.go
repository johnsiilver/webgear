package html

import (
	"context"
	"strings"
	"testing"
)

func TestTextArea(t *testing.T) {
	tests := []struct {
		desc string
		textArea *TextArea
		want string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			textArea: &TextArea{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Name: "name",
				Form: "formid",
				Cols: 10,
				MaxLength: 1000,
				Rows: 40,
				DirName: "name.dir",
				Wrap: HardWrap,
				Placeholder: "placeholder",
				AutoFocus: true,
				Disabled: true,
				ReadOnly: true, 
				Required: true,
				Element: TextElement("text"),
				Events:   (&Events{}).AddScript(OnError, "handleError"),
			},

			want: strings.TrimSpace(`
<textarea accesskey="key" onerror="handleError" name="name" form="formid" cols="10" maxlength="1000" rows="40" dirname="name.dir" wrap="hard" placeholder="placeholder" autofocus disabled readonly required>text
</textarea>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.textArea.Execute(pipe)
		if test.want != got.String() {
			t.Errorf("TestTextArea(%s): \n\tgot  %q\n\twant %q", test.desc, got, test.want)
		}
	}
}
