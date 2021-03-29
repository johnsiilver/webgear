<p align="right">
  ⭐ &nbsp;&nbsp;<strong>the project to show your appreciation.</strong> :arrow_upper_right:
</p>

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)](https://pkg.go.dev/github.com/johnsiilver/webgear)

WebGear provides libraries to make Go the single language for writing web applications.

## OK, Why would I use this?

Good question. The world of web development has a lot of choices today. You could use React, Angular, Vuze, ...

On top of that, this toolkit is written by someone who loathes web development. It may be ubiquitous, but HTTP
and HTML (in its modern form) are absolutely the worst. Javascript? are you kidding me with that language??
Typescript, a slightly less bad Javascript that has to compile to Javacript? Dart????

So, if you feel similar to that last paragraph, this might be for you:
- Program in Go and serve directly from Go.
- No more template debugging
- Serve your assets easily

You can still insert Javascript if you need to. But you can also go to WASM and leave Javascript behind.

You still must deal with CSS, but I provide web components so you can keep your CSS sane.

## Production Quality?

No idea. I've used it for personal projects and such. It has some rough edges. I don't support all tags (there are a lot).
Its performant for what I do, but if I was Google's front page I doubt it.

Also, its going to change. I don't expect any major "destroy everything" type of changes. But I'm sure it will get lots
of refinements. Most will go unnoticed. 

Also, I doubt I will ever release a 1.0. I'm not in love with Go's semantic versioning where I have to create new v2/ directories
and such. I get why, I just don't like it.

## Is this a Framework?

Not really, though I imagine you could build a framework from this. All of the code you will find here are at the low 
HTML and Javascript level. 

It might be argued that the addition of the Component type and some WASM helpers are frameworky, but I included them
to simply get around the nastiness of using some very useful web concepts.

## Packages

There are several main packages:
- html/ - Provides HTML tags as Go types
- html/builder - Allows dynamic building of HTML documents
- handlers/ - Provides http.Handle(s) that serve content and files
- wasm/ - Provides tooling to build WASM apps wihtout interacting with syscall/js

More indepth documentation will be in the godoc.

## Examples

### Example: Create a simple HTML page and serve it

This can be found in html/examples/basic .

A few notes:
- A "*Doc" object represents an HTML document and is the fundamental structure used.
- An "Element" interface is used to represent all HTML elements that are defined.
- I use the "." import to avoid having to type &html.Div{} and &html.Title{} (all Go rules have exceptions, I think this is one)


```go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/johnsiilver/webgear/handlers"

	. "github.com/johnsiilver/webgear/html"
)

var (
	dev = flag.Bool("dev", false, "Prevents the browser from caching content when doing development")
	port = flag.Int("port", 9568, "The port to server on")
)

func main() {
	flag.Parse()

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "utf-8"},
				// TextElement() creates an Element that contains the text passed.
				&Title{TagValue: TextElement("Hello World")},
				&Link{
					Rel: "stylesheet",
					// URLParse() takes a string and create a *url.URL object
					// or panics.
					Href: URLParse("/static/index/index.css"),
				},
			},
		},
		Body: &Body{
			Elements: []Element{
				&H{
					GlobalAttrs: GlobalAttrs{
						Class: "pageText",
					},
					Level: 1,
					Elements: []Element{
						TextElement("Hello World"),
					},
				},
			},
		},
	}

	// Preventing the browser from caching is helpful during development.
	opts := []handlers.Option{}
	if *dev {
		opts = append(
			opts,
			handlers.DoNotCache(),
		)
	}

	// Create our http.Handler. By default we use level 4 gzip compression.
	h := handlers.New(opts...)

	// Serves up files ending with .css from /static/...
	// NOTE: This will not serve any content from the root directory, all files to be included
	// must live in sub-directories.
	// NOTE: This will never serve files with .go and a few other extensions.
	h.ServeFilesWorkingDir([]string{".css"})

	// Our doc will now be served at the index page.
	h.MustHandle("/", doc)

	// Serve the content using the http.Server.
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        h.ServerMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("http server serving on :%d", *port)
	log.Fatal(server.ListenAndServe())
}```

This is the most basic way of creating a page. If writing out a lot of content, you run into a flaw with the example above::

- The Doc becomes this nasty indented mountain going off the right of your screen, just like HTML

For that kind of problem or generating some types of dynamic content, I suggest the builder/ package.

### Example: Use the builder to build a sample HTML page

The html/builder/ package allows building pages in a way that is more dynamic and prevents the
code from going off the right of your screen. 

In the example below, we won't have all the boilerplate as we did before.  

Also, while we are statically building a table here, this really shines when building things with for loops.

The complete code can be found in html/examples/builder .

```go

...

// Note: You can also build the head with the builder package.
head := Head{
	Elements: []Element{
		&Meta{Charset: "utf-8"},
		&Title{TagValue: TextElement("Hello World")},
		&Link{
			Rel: "stylesheet",
			Href: URLParse("/static/index/index.css"),
		},
	},
}

// Create a builder.HTML object with an initial head and body.
build := builder.NewHTML(&head, &Body{})

// Add a single H1 element containing some text.
build.Add(
	&H{Level: 1, Elements: []Element{TextElement("Movies I Like")}},
)

// Create a div with a set ID and move the builders context inside the div.
build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "myTableDiv"}})
// Create a table in the div and move the builders context inside the table.
build.Into(&Table{})
// Add a tr to the table and move into it.
build.Into(&TR{})
// Add two th elements to the table.
build.Add(
	&TH{Element: TextElement("Movie")},
	&TH{Element: TextElement("Category")},
)
// Use Up() to set the context back to the table.
build.Up() // We are now inside the table
build.Into(&TR{})
build.Add(
	&TD{Element: TextElement("Blade Runner")},
	&TD{Element: TextElement("SciFi")},
)
build.Up() // We are now inside the table
build.Into(&TR{})
build.Add(
	&TD{Element: TextElement("Memento")},
	&TD{Element: TextElement("Drama")},
)
build.Up().Up() // We are now inside the div

...

// Our doc will now be served at the index page.
h.MustHandle("/", build.Doc())

...

```

### Example: Use the dynamic type to respond with dynamic content

When doing static HTML (not WASM, which can respond to dynamic needs in a cleaner way)), 
you often need to build your HTML on the server based on things like URL paths, query strings, or the POST body.

Enter the Dynamic() and DynamicFunc() type.

Simply put, you can add a Dynamic() to create an Element that uses a DynamicFunc type when the doc is executed.
It will pass a Pipeline object to your DynamicFunc() which responds with the []Element you want to add in that location.
The Pipeline object will contain your http.Request object.

```go
// HelloUser looks for the user's name as a query string element and prints hello to that name.
func HelloUser(pipe Pipeline) []Element {
	name := pipe.Req.URL.Query().Get("name")
	if name == "" {
		return []Element{&H{Level: 2, Elements: []Element{TextElement("Hello Unknown User")}}}
	}

	return []Element{&H{Level: 2, Elements: []Element{TextElement(fmt.Sprintf("Hello %s", name))}}}
}

func main() {
	flag.Parse()

	doc := &Doc{
		Head: &Head{
			Elements: []Element{
				&Meta{Charset: "utf-8"},
				&Title{TagValue: TextElement("Hello Person")},
			},
		},
		Body: &Body{
			Elements: []Element{Dynamic(HelloUser)},
		},
	}

	...
}

```

This example will simply read the query string and if there is a "name" key will print hello to that name.

But Dynamic can do much more. You can create a type that has access to databases or clients to other services
and then use a method on that type to implement the DynamicFunc(). Then that func can grab all kinds of data or record data
in response to a request.

I find for many of my needs I never need to go to WASM or Javascript using just Dynamic.

### Example: What about Events?

Let's create an event that makes a modal and when a button is pressed it hides. 

This will use a little javascript, as javascript is needed with non-WASM events.

```go
...

func modal() []Element {
        build := builder.NewHTML(&Head{}, &Body{})

        build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "container", Class: "container"}})
        build.Into(&Div{GlobalAttrs: GlobalAttrs{Class: "card"}})
        build.Into(&Div{GlobalAttrs: GlobalAttrs{Class: "content-wrapper"}})
        build.Add(&P{Elements: []Element{TextElement("Demo Modal")}})
        build.Up()
        build.Add(
                &Span{
                        GlobalAttrs: GlobalAttrs{Class: "button button__link"},
                        Elements: []Element{TextElement("Close")},
                        Events: (&Events{}).AddScript(
                                OnClick,
                                `document.getElementById('container').style.visibility = 'hidden';`,
                        ),
                },
        )

        return build.Doc().Body.Elements
}

func main() {
	flag.Parse()

	doc := &Doc{
                Head: &Head{
                        Elements: []Element{
                                &Meta{Charset: "utf-8"},
                                &Title{TagValue: TextElement("Hello World")},
                                &Link{
                                        Rel: "stylesheet",
                                        Href: URLParse("/static/index/index.css"),
                                },
                        },
                },
                Body: &Body{
                        Elements: modal(),
                },
        }

	...

}

```

### Example: Use web components for cleaner CSS

Ok, if you've done HTML/CSS for any length of time you hate when you have different elements of the page which are
affecting by other element's CSS styling.

Nothing is worse than changing a CSS element and watching the whole page change.

Some years ago someonecame up with the idea of web components:
https://developer.mozilla.org/en-US/docs/Web/Web_Components

My first experience with these was in Dart using Polymer. However after they canned the Dart version, I decided I never wanted to use
either tech again (nope, not even with Flutter).

But having reusable components (though I think this is less compelling) and having custom tags with encapsulated CSS 
was somethign I missed.

However adding components way more difficult than it had to be, relying on Javascript nastiness.

So I have created a Component() Element to allow you to create these encapsulated elements with ease.

In the example below, I'm going to create a component that acts as a banner. The code is located at html/examples/components .

Defining the component:
```go
package banner

import (
        "github.com/johnsiilver/webgear/component"
        "github.com/johnsiilver/webgear/html/builder"

        . "github.com/johnsiilver/webgear/html"
)

// New constructs a new component that shows a banner.  
func New(name string, options ...component.Option) (*component.Gear, error) {
        build := builder.NewHTML(&Head{}, &Body{GlobalAttrs: GlobalAttrs{ID: "banner"}})
        build.Into(&Div{GlobalAttrs: GlobalAttrs{ID: "banner"}})
        build.Add(&Link{Rel: "stylesheet", Href: URLParse("/static/banner/banner.css")})

        build.Into(&A{Href: URLParse("/")})
        build.Add(&Img{GlobalAttrs: GlobalAttrs{ID: "gopher"}, Src: URLParse("/static/banner/scientist.svg")})
        build.Up()

        build.Into(&A{Href: URLParse("/")})
        build.Into(&Span{GlobalAttrs: GlobalAttrs{ID: "title"}})
        build.Add(TextElement("Example Banner"))

        gear, err := component.New(name, build.Doc(), options...)
        if err != nil {
                return nil, err
        }

        return gear, nil
}
```

Create a page that uses the component:
```go
package index

import (
        "github.com/johnsiilver/webgear/html/examples/components/banner"

        . "github.com/johnsiilver/webgear/html"
)

const (
        bannerGearName  = "banner-component"
)

// New creates a new Doc object that can renders the index page.
func New() (*Doc, error) {
	// Create our banner object. The name passed must be two parts with a "-" between them (part of the web component standard).
	// We don't create a name because you can spawn a gear multiple times that act differently and are added with their own tags.
        bannerGear, err := banner.New(bannerGearName)
        if err != nil {
                return nil, err
        }

        doc := &Doc{
                Head: &Head{
                        Elements: []Element{
                                &Meta{Charset: "UTF-8"},
                                &Title{TagValue: TextElement("Go Language Basics")},
                                &Link{Rel: "stylesheet", Href: URLParse("/static/index/index.css")},
                                &Link{Href: URLParse("https://fonts.googleapis.com/css2?family=Share+Tech+Mono&display=swap"), Rel: "stylesheet"},
                                &Link{Href: URLParse("https://fonts.googleapis.com/css2?family=Nanum+Gothic&display=swap"), Rel: "stylesheet"},
                        },
                },
                Body: &Body{
                        Elements: []Element{
				// Okay, YOU MUST put the Gear here before you do the component. That causes all the necessary
				// javascript stuff to happen.
                                bannerGear, // This causes the code to render.

				// This causes the component to be inserted in this particular place in the Body.
				// Components are their own tags, so you will see in the HTML: "<banner-component id="banner"></banner-component>"
                                &Component{GlobalAttrs: GlobalAttrs{ID: "banner"}, Gear: bannerGear},
                        },
                },
        }

        return doc, nil
}
```

Now you just need to serve it:

```go
package main

...

func main() {
	flag.Parse()

	doc, err := index.New()
        if err != nil {
                panic(err)
        }

	...

	 // Serves up files ending with .css from /static/...
        h.ServeFilesWorkingDir([]string{".svg", ".css"})

        // Our doc will now be served at the index page.
        h.MustHandle("/", doc)

	...
}
```

### Example: Screw Javascript, do some WASM!!

I'm going to add examples on how to use the WASM framework. Right now it works, but I haven't used it for a lot of
things. You are free to use it, but I'm going to need some more time to document it.

## Compile and Runtime errors

This package has **some** protections from bad code. The problem with HTML is that what was illegal today is not
always going to be the case. 

There is some compile level checking done with types. There are also some minor validations with the ability to add more.

But there are also runtime panics in the builder, as it heavily uses Element and not all Elements can be added to
all elements. These panics either happen during Doc creation (so no big deal) or in isolated sections and the web server
will just generate a 404. In practice I have not found this to be an issue with uptime.

Here's an example:

```
panic: reflect.Set: value of type *html.TD is not assignable to type html.TableElement

goroutine 1 [running]:
reflect.Value.assignTo(0x137f7c0, 0xc00020c2d0, 0x16, 0x13d2cab, 0xb, 0x137dae0, 0xc00014dda0, 0x97, 0xc00012b7e0, 0x1356d20)
	/usr/local/go/src/reflect/value.go:2425 +0x405
reflect.Value.Set(0x137dae0, 0xc00014dda0, 0x194, 0x137f7c0, 0xc00020c2d0, 0x16)
	/usr/local/go/src/reflect/value.go:1554 +0xbd
reflect.Append(0x1356d20, 0xc0002071b8, 0x197, 0xc00017dd68, 0x1, 0x1, 0xc0002071b8, 0x197, 0xc00020c3c0)
	/usr/local/go/src/reflect/value.go:2037 +0xea
github.com/johnsiilver/webgear/html/builder.(*HTML).Add(0xc00017de50, 0xc00017de90, 0x2, 0x2)
	/Users/jdoak/trees/webgear/html/builder/builder.go:111 +0x39b
main.main()
	/Users/jdoak/trees/webgear/html/examples/builder/builder.go:63 +0x7d2
```

## FAQ

## I have a panic with a backtrace that is hard to track down

It is possible you will see something like:
```
2021/03/28 15:19:20 template: doc:5:8: executing "doc" at <.Self.Body.Execute>: error calling Execute: template: body:4:3: executing "body" at <.Execute>: error calling Execute: runtime error: invalid memory address or nil pointer dereference
```

Yuck, I hear ya. This usually means a bug between the object and the template or you forgot to put something vital in an element and I forgot
to validate that the vital part is there.

You can open a bug request and I'll need to have a look at the code to try and track it down.

## Sometimes it looks like a page is getting called multiple times

This is almost always because some page is being called that doesn't exist and is redirected to /.
Most of the time that is the favicon. Add a favicon and that will probably stop.

## Do I have to recompile on every change?

The answer is maybe. CSS when passing the NoCache option will always load on every page refresh. However,
changes to the Go code require a recompile.

Fixing CSS on the fly is convenient, which is why I did not try and implement CSS in Go.

## Not every tag, every tag option, or every event is in here

Well, that's true. HTML5 is a humongous spec. Frankly, it is needlessly big. I'd say about 70% of it could be
thrown in a trash bin and everything could still be done.

There are a couple of things that you can do if you run into this:
- Put in a feature request
- Look at the existing code, add it and send it to me
- Try using a different tag that will accomplish the same thing

## Acknowledgements

This project is built on top of the work of all the Go Authors of course. But in this case I'd like to make a special call out
to Richard Musiol, author of GopherJS and syscall/js. That is a crazy amount of work to try and get away from Javascript and of course the WASM
part of this would be nowhere without him.
