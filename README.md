# Regression Wasm

This repository serves a simple [web assembly](https://webassembly.org/) (wasm) application 
to perform a regression, using data from a table in the browser, which can be loaded as a delimited file
by the user. We use a simple [regression library](https://github.com/sajari/regression) to do
the work.

## About

### Why?

Web assembly can allow us to interact with compiled code directly in the browser,
doing away with any need for a server. While I don't do a large amount of data analysis
for my role proper, I realize that many researchers do, and so with this in mind, 
I wanted to create a starting point for developers to interact with data in the browser.
The minimum conditions for success meant:

 1. being able to load a delimited file into the browser
 2. having the file render as a table
 3. having the data be processed by a compiled wasm
 4. updating a plot based on output from 3.

Thus, the application performs a simple regression based on loading data in the table,
and then plotting the result. To make it fun, I added a cute gopher logo and used an xkcd
plotting library for the result.

### Customization

The basics are here for a developer to create (some GoLang based) functions to
perform data analysis on an input file, and render back to the screen as a plot.
If you need any help, or want to request a custom tool, please don't hesitate to
[open up an issue](https://www.github.com/vsoch/regression-wasm/issues).

## Development

### Local

If you are comfortable with GoLang, and have installed [emscripten](https://emscripten.org), 
you can clone the repository into your $GOPATH under the github folder:

```bash
$ mkdir -p $GOPATH/src.github.com/vsoch
$ cd $GOPATH/src.github.com/vsoch
$ git clone https://www.github.com/vsoch/regression-wasm
```

And then build the wasm.

```bash
$ make
```

And cd into the "docs" folder and start a server to see the result.

```bash
$ cd docs
$ python -m http.server 9999
```

Open the browser to http://localhost:9999


## Docker

If you don't want to install dependencies, just clone the repository, and
build the Docker image:

```bash
$ docker build -t vanessa/regression-wasm .
```

It will install [emscripten](https://emscripten.org/docs/getting_started/FAQ.html),
add the source code to the repository, and compile to wasm. You can then
run the container and expose port 80 to see the compiled interface:

```bash
$ docker run -it --rm -p 80:80 vanessa/regression-wasm
``` 

Then you can proceed to use the interface.
