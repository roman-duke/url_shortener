package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[35m"
	Purple = "\033[36m"
	White  = "\033[37m"
)

func main() {
	component := app()

	// ============================================= //
	// =========== Serve static assets ============= //
	dir := http.Dir("./cmd/web/assets")
	fS := http.FileServer(dir)

	handler := http.StripPrefix("/static", fS)

	http.Handle("/static/", handler)
	// ============================================= //

	http.Handle("/", templ.Handler(component))

	fmt.Printf("%sServer listening at http://localhost:8080\n%s", Cyan, Reset)

	http.ListenAndServe(":8080", nil)
}
