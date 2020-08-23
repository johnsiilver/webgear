package builder

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"

	. "github.com/johnsiilver/webgear/html"
)

func TestHTML(t *testing.T) {
	b := NewHTML(&Head{}, &Body{})

	b.Into(&Div{})      // Adds the div and moves into the div
	b.Into(&Table{})    // Adds the table and moves into the table
	b.Into(&TR{})       // Adds the table row and moves into the table row
	b.Add(&TD{}, &TD{}) // Adds two table role elements, but stays in the row.
	b.Up().Up()         // We now move back to the table, if we called b.Up() again, we'd be at the div
	b.Into(&Table{})    // Adds the table and moves into the table
	b.Into(&TR{})       // Adds the table row and moves into the table row
	b.Add(&TD{}, &TD{})

	want := &Doc{
		Head: &Head{},
		Body: &Body{
			Elements: []Element{
				&Div{
					Elements: []Element{
						&Table{
							Elements: []TableElement{
								&TR{
									Elements: []TRElement{
										&TD{},
										&TD{},
									},
								},
							},
						},
						&Table{
							Elements: []TableElement{
								&TR{
									Elements: []TRElement{
										&TD{},
										&TD{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if diff := pretty.Compare(want, b.Doc()); diff != "" {
		t.Errorf("TestHTML: -want/+got:\n%s", diff)
	}
}
