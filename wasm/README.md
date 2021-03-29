NOTE: This is a work in progress and neither the doc nor the package is ready for rock and roll.

WebGear WASM provides packages to simplify using WASM in Go utilizing by boostrapping all files needed for WASM, utilizing the Webear html package and providing UI DOM manipulation without having to get into syscall/JS.

## Why Use This?

- You dislike Javascript
- You like the other WebGear packages
- You don't want to figure out how to bootstap WASM

## Production Ready?

Nope. 

Look, WASM is still somewhat experimental in Go. WASM itself is
designed for C++ and doesn't have GC hooks. So binaries are bigger
than they should be. Yeah, there is tinyGo, but that starts getting
way harder.

Here's the thing: is your app a large codebase that is revenue generating at scale?  If so, go back to JS. 

If not (99% of everyone out there), you can try something different.

## Concepts

### How WASM bootstrapping works

Conceptually, WASM is like Javascript. With JS you push code from the server to the browswer and the browser runs the code.

The big differences are that with WASM, you push Javascript and HTML to the browser that then fetches your WASM from the server and then runs it.

The Javascript for bootstrapping is provided by the Go team. But you have to serve it and setup the WASM download.

Luckily for you we automate that whole thing so you don't have to worry
about it.

### WASM code is compiled for a different environment

When you write your WASM code, you will be compiling for a new target. You can't run tests like normal because you can't run WASM in your environment. Basically any package that uses syscall/js cannot be tested in your normal environment.

### When you want to see changes, you have to recompile the WASM and re-run your server

You can't just "go run main.go". Your server is serving a WASM file,
 that WASM file is its own compiled binary.

As you will see below, this just means there are two compilations that
must happen.


