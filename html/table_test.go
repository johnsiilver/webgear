package html

import (
	"context"
	"strings"
	"testing"
)

func TestTable(t *testing.T) {
	tests := []struct {
		desc  string
		table *Table
		want  string
	}{
		{
			desc: "All attributes + 1 global + 1 event ",
			table: &Table{
				GlobalAttrs: GlobalAttrs{
					AccessKey: "key",
				},
				Events: (&Events{}).OnError("handleError"),
				Elements: []TableElement{
					&Caption{
						GlobalAttrs: GlobalAttrs{
							AccessKey: "key",
						},
						Events:  (&Events{}).OnError("handleError"),
						Element: TextElement("caption"),
					},
					&ColGroup{
						GlobalAttrs: GlobalAttrs{
							AccessKey: "key",
						},
						Events: (&Events{}).OnError("handleError"),
						Span:   1,
						Elements: []ColGroupElement{
							&Col{
								GlobalAttrs: GlobalAttrs{
									AccessKey: "key",
								},
								Events: (&Events{}).OnError("handleError"),
							},
							&Col{
								GlobalAttrs: GlobalAttrs{
									AccessKey: "key",
								},
								Events: (&Events{}).OnError("handleError"),
								Span:   1,
							},
						},
					},
					&THead{
						GlobalAttrs: GlobalAttrs{
							AccessKey: "key",
						},
						Events: (&Events{}).OnError("handleError"),
						Elements: []*TR{
							&TR{
								Elements: []TRElement{
									&TH{
										GlobalAttrs: GlobalAttrs{
											AccessKey: "key",
										},
										Events:  (&Events{}).OnError("handleError"),
										Element: TextElement("col 0 header"),
									},
									&TH{
										GlobalAttrs: GlobalAttrs{
											AccessKey: "key",
										},
										Events:  (&Events{}).OnError("handleError"),
										Element: TextElement("col 1 header"),
									},
								},
							},
						},
					},
					&TBody{
						GlobalAttrs: GlobalAttrs{
							AccessKey: "key",
						},
						Events: (&Events{}).OnError("handleError"),
						Elements: []*TR{
							&TR{
								Elements: []TRElement{
									&TD{
										Element: TextElement("row0col0"),
									},
									&TD{
										Element: TextElement("row0col1"),
									},
								},
							},
							&TR{
								Elements: []TRElement{
									&TD{
										Element: TextElement("row1col0"),
									},
									&TD{
										Element: TextElement("row1col1"),
									},
								},
							},
						},
					},
					&TFoot{
						GlobalAttrs: GlobalAttrs{
							AccessKey: "key",
						},
						Events: (&Events{}).OnError("handleError"),
						Elements: []*TR{
							&TR{
								Elements: []TRElement{
									&TD{
										Element: TextElement("footer0"),
									},
									&TD{
										Element: TextElement("footer1"),
									},
								},
							},
						},
					},
				},
			},

			want: strings.TrimSpace(`
<table accesskey="key" onerror="handleError">
	<caption accesskey="key" onerror="handleError">
		caption
	</caption>
	<colgroup span="1" accesskey="key" onerror="handleError">
		<col accesskey="key" onerror="handleError" >
		<col span="1" accesskey="key" onerror="handleError" >
	</colgroup>
	<thead accesskey="key" onerror="handleError">
		<tr>
			<th element="col 0 header" accesskey="key" onerror="handleError">
				col 0 header
			</th>
			<th element="col 1 header" accesskey="key" onerror="handleError">
				col 1 header
			</th>
		</tr>
	</thead>
	<tbody accesskey="key" onerror="handleError">
		<tr>
			<td>
				row0col0
			</td>
			<td>
				row0col1
			</td>
		</tr>
		<tr>
			<td>
				row1col0
			</td>
			<td>
				row1col1
			</td>
		</tr>
	</tbody>
	<tfoot accesskey="key" onerror="handleError">
		<tr>
			<td>
				footer0
			</td>
			<td>
				footer1
			</td>
		</tr>
	</tfoot>
</table>
`),
		},
	}

	for _, test := range tests {
		got := &strings.Builder{}
		pipe := NewPipeline(context.Background(), nil, got)
		test.table.Execute(pipe)
		if removeSpace(test.want) != removeSpace(got.String()) {
			t.Errorf("TestTable(%s): \n\tgot  %q\n\twant %q", test.desc, removeSpace(got.String()), removeSpace(test.want))
		}
	}
}
