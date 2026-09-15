package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/roman-duke/url_shortener/cmd/web/components"
	"github.com/roman-duke/url_shortener/shortener"
	"github.com/roman-duke/url_shortener/storage/memory"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

var service shortener.Service
var store memory.MyStore

func shortenHandler(val string) templ.Component {
	// call the service layer to create a corresponding short link
	shortLink, err := service.Shorten(val)

	return components.ShortenComponent(val, shortLink.Code, err)
}

func resolveHandler(code string) templ.Component {
	// call the service layer to resolve the short code
	sLink, err := service.Resolve(code)

	return components.ResolveComponent(sLink.LongUrl, err)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/shorten":
		// Extract the link field from the form
		// and call the respective service layer
		l := r.FormValue("link")
		component := shortenHandler(l)

		// return the rendered html
		component.Render(r.Context(), w)

	case "/resolve":
		// Extract just the code field from the
		// form and call the respective service layer
		c := r.FormValue("code")
		component := resolveHandler(c)

		// return the rendered html
		component.Render(r.Context(), w)
	}
}

// func getHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	component := app()

	// ============================================= //
	// =========== Serve static assets ============= //
	dir := http.Dir("./cmd/web/assets")
	fS := http.FileServer(dir)

	handler := http.StripPrefix("/static", fS)

	http.Handle("/static/", handler)
	// ============================================= //

	// =========== Create the service layer ======== //
	service = *shortener.NewService(&store)
	// ============================================= //

	http.Handle("/", templ.Handler(component))
	http.HandleFunc("POST /shorten", postHandler)
	http.HandleFunc("POST /resolve", postHandler)

	fmt.Printf("%sServer listening at http://localhost:8080\n%s", Cyan, Reset)

	http.ListenAndServe(":8080", nil)
}
