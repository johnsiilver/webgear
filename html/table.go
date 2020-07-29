package html

import (
	"html/template"
	"strings"
)

var tableTmpl = template.Must(template.New("table").Parse(strings.TrimSpace(`
<table {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</table>
`)))

// TableElement represents a tag that can be contained inside a Table.
type TableElement interface {
	Element
	isTableElement()
}

// Table represents a division tag.
type Table struct {
	GlobalAttrs
	// Elements are elements contained within the Table.
	Elements []TableElement

	Events *Events
}

func (t *Table) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := tableTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var trTmpl = template.Must(template.New("tr").Parse(strings.TrimSpace(`
<tr {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</tr>
`)))

// TRElement is an element that can be included in a TR tag.
type TRElement interface {
	Element
	isTRElement()
}

// TR is a table row for use inside a Table tag.
type TR struct {
	GlobalAttrs
	Events *Events

	Elements []TRElement
}

func (t TR) isTableElement() {}

func (t *TR) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := trTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

// THScope is the scope of various table elements.
type THScope int

const (
	ColScope      THScope = 1
	ColGroupScope THScope = 2
	RowScope      THScope = 3
	RowGroupScope THScope = 4
)

var thTmpl = template.Must(template.New("th").Parse(strings.TrimSpace(`
<th {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{.Self.Element.Execute $data}}
</th>
`)))

// TH is a table header for use inside a Table tag.
type TH struct {
	GlobalAttrs
	Events *Events

	// Element contains the text in the TH tag.
	Element TextElement

	// Abbr specifies an abbreviated version of the content in a header cell.
	Abbr string
	// ColSpan specifies the number of columns a header cell should span.
	ColSpan int
	// Headers specifies one or more header cells a cell is related to.
	Headers []string
	// RowSpan specifies the number of rows a header cell should span.
	RowSpan int
	// Scope specifies whether a header cell is a header for a column, row, or group of columns or rows.
	Scope THScope
}

func (t TH) isTRElement() {}

func (t *TH) Attr() template.HTMLAttr {
	output := structToString(t)
	return template.HTMLAttr(output)
}

func (t *TH) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := thTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var tdTmpl = template.Must(template.New("td").Parse(strings.TrimSpace(`
<td {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{.Self.Element.Execute $data}}
</td>
`)))

// TD holds table data inside a TR tag inside a Table tag.
type TD struct {
	GlobalAttrs
	Events *Events

	Element Element

	// ColSpan specifies the number of columns a header cell should span.
	ColSpan int
	// Headers specifies one or more header cells a cell is related to.
	Headers []string
	// RowSpan specifies the number of rows a header cell should span.
	RowSpan int
}

func (t TD) isTRElement() {}

func (t *TD) Attr() template.HTMLAttr {
	output := structToString(t)
	return template.HTMLAttr(output)
}

func (t *TD) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := tdTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var captionTmpl = template.Must(template.New("caption").Parse(strings.TrimSpace(`
<caption {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{.Self.Element.Execute $data}}
</caption>
`)))

// Caption is a table element for use inside a Table tag.
type Caption struct {
	GlobalAttrs
	Events *Events

	// Element contains the text in the TH tag.
	Element TextElement
}

func (*Caption) isTableElement() {}

func (c *Caption) Execute(pipe Pipeline) string {
	pipe.Self = c

	if err := captionTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var colGroupTmpl = template.Must(template.New("colGroup").Parse(strings.TrimSpace(`
<colgroup {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</colgroup>
`)))

// ColGroupElement is a tag that can be included in a ColGroup.
type ColGroupElement interface {
	Element
	isColGroupElement()
}

// ColGroup is a table tag that specifies one or more columns in a table for formatting. It is a applied to a Table
// element after caption tags and before THead, TBody and TFoot.
type ColGroup struct {
	GlobalAttrs
	Events *Events

	// Elements contains the Col tags associated with this ColGroup.
	Elements []ColGroupElement

	// Span specifies the number of columns a header cell should span.
	Span int
}

func (*ColGroup) isTableElement() {}

func (c *ColGroup) Attr() template.HTMLAttr {
	output := structToString(c)
	return template.HTMLAttr(output)
}

func (c *ColGroup) Execute(pipe Pipeline) string {
	pipe.Self = c

	if err := colGroupTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var colTmpl = template.Must(template.New("col").Parse(strings.TrimSpace(`
<col {{.Self.Attr}} {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
`)))

// Col is a tag that specifies column properties within a ColGroup.
type Col struct {
	GlobalAttrs
	Events *Events

	// Span specifies the number of columns a header cell should span.
	Span int
}

func (*Col) isColGroupElement() {}

func (c *Col) Attr() template.HTMLAttr {
	output := structToString(c)
	return template.HTMLAttr(output)
}

func (c *Col) Execute(pipe Pipeline) string {
	pipe.Self = c

	if err := colTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var tHeadTmpl = template.Must(template.New("tHead").Parse(strings.TrimSpace(`
<thead {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</thead>
`)))

// THead is used to group header content in a table.
type THead struct {
	GlobalAttrs
	Events *Events

	Elements []*TR
}

func (t *THead) isTableElement() {}

func (t *THead) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := tHeadTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var tBodyTmpl = template.Must(template.New("tBody").Parse(strings.TrimSpace(`
<tbody {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</tbody>
`)))

// TBody is used to group the body elements in an Table.
type TBody struct {
	GlobalAttrs
	Events *Events

	Elements []*TR
}

func (t *TBody) isTableElement() {}

func (t *TBody) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := tBodyTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}

var tFootTmpl = template.Must(template.New("tFoot").Parse(strings.TrimSpace(`
<tfoot {{.Self.GlobalAttrs.Attr}} {{.Self.Events.Attr}}>
	{{- $data := .}}
	{{- range .Self.Elements}}
	{{.Execute $data}}
	{{- end}}
</tfoot>
`)))

// THead is used to group footer content in a table.
type TFoot struct {
	GlobalAttrs
	Events *Events

	Elements []*TR
}

func (t *TFoot) isTableElement() {}

func (t *TFoot) Execute(pipe Pipeline) string {
	pipe.Self = t

	if err := tFootTmpl.Execute(pipe.W, pipe); err != nil {
		panic(err)
	}

	return EmptyString
}
