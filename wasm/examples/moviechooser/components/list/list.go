package list

import (
	"log"
	"sort"
	"strings"
	"syscall/js"

	"github.com/johnsiilver/webgear/component"
	"github.com/johnsiilver/webgear/html/builder"
	"github.com/johnsiilver/webgear/wasm"

	. "github.com/johnsiilver/webgear/html"
)

type Filters []string

type Movie struct {
	Name       string
	Categories Categories
}

type Categories map[string]bool

func (c Categories) Add(s string) Categories {
	c[s] = true
	return c
}
func (c Categories) Has(category string) bool {
	_, ok := c[category]
	return ok
}
func (c Categories) String() string {
	s := []string{}
	for k := range c {
		s = append(s, k)
	}
	sort.Strings(s)
	return strings.Join(s, ", ")
}

var CategoriesList = []string{
	Action,
	Drama,
	Comedy,
	Romance,
	SciFi,
	Fantasy,
	Documentary,
}

const (
	Action      = "Action"
	Drama       = "Drama"
	Comedy      = "Comedy"
	Romance     = "Romance"
	SciFi       = "SciFi"
	Fantasy     = "Fantasy"
	Documentary = "Documentary"
)

var CategoryList = []string{
	Action,
	Drama,
	Comedy,
	Romance,
	SciFi,
	Fantasy,
	Documentary,
}

var Movies = []Movie{
	{
		"Blade Runner",
		Categories{}.Add(Action).Add(Drama).Add(Romance).Add(SciFi),
	},
	{
		"Apocolypse Now",
		Categories{}.Add(Action).Add(Drama),
	},
	{
		"Ace Ventura, Pet Detective",
		Categories{}.Add(Action).Add(Comedy).Add(Fantasy),
	},
	{
		"Love Actually",
		Categories{}.Add(Action).Add(Comedy).Add(Drama).Add(Romance),
	},
	{
		"Enron, The Smarest Men in the Room",
		Categories{}.Add(Documentary),
	},
	{
		"Star Trek: The Wrath of Khan",
		Categories{}.Add(Action).Add(SciFi).Add(Drama),
	},
}

// New constructs a new component that shows a list of categories that are selectable.
func New(name string, filters Filters, options ...component.Option) (*component.Gear, error) {
	b := builder.NewHTML(&Head{}, &Body{})
	b.Add(&Link{Rel: "stylesheet", Href: URLParse("/static/components/list/list.css")})

	b.Into(&Table{GlobalAttrs: GlobalAttrs{ID: "movieList"}})
	// Add our table header.
	b.Into(&THead{})
	b.Into(&TR{})
	b.Add(&TH{Element: TextElement("Name")})
	b.Add(&TH{Element: TextElement("Categories")})
	b.Up().Up()
	// Add our table body.
	b.Into(&TBody{})
	// This is a bad data structure.  I would not do this in real life, but
	// this is an example so I didn't want to write something complicated.  Also, this
	// would be done on the backend and not a list in the client code.
	for _, movie := range Movies {
		ok := true
		if len(filters) != 0 {
			for _, filter := range filters {
				if !movie.Categories.Has(filter) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
		}
		b.Into(
			wasm.AttachListener(
				LTClick,
				false,
				func(this, root js.Value) {
					log.Println("hello")
				},
				&TR{GlobalAttrs: GlobalAttrs{ID: movie.Name}},
			),
		)
		b.Add(&TD{Element: TextElement(movie.Name)})
		b.Add(&TD{Element: TextElement(movie.Categories.String())})
		b.Up()
	}

	gear, err := component.New(name, b.Doc(), options...)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
