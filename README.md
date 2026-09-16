# URL SHORTENER (Single-node version)
Simple URL Shortener project written in Go. The underlying architecture of this project
is the [hexagonal architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)).
The domain (the core url shortening service) sits as a standalone entity, then drives and is driven
by other adapters.

## Driving Adapters (who drives the app - e.g. web, CLI)
The main entry points (the channels which can communicate with the application) include:
- the Web layer (Go + htmx + templ)
- the CLI.

## Driven Adapters (what the app drives - e.g. db, internal memory, etc)
In this application, the only driven adapter that exists is the persistence layer. For the
sake of simiplicity, there is only one - "internal memory".

### Running the code
A makefile has been included to orchestrate necessasry task runners required for running the
web-layer of this project. Simply run the following from the root.

```
make dev
```

And then for necessary cleanup (given that it installs the tailwindcss binary and runs it as a
watcher amongst other things), run the following:

```
make clean
```

For the CLI-layer, it can easily be run as follows:

``` go
go run ./cmd/cli
```


> _Perseverance is not a long race; it is many short races one after the other._ ~ **Walter Elliot**
