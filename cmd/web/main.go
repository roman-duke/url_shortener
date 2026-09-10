package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

// We generally should have a constructor function (idiomatic Go)

// this is cool for generating static data
// as you can see, any argument passed to the templ component is rendered
// at runtime and we have no way to generate dynamic data.

func main() {
	http.Handle("/", templ.Handler(appShell("Eren Jaeger", time.Now())))
	http.Handle("/404", templ.Handler(notFoundComponent()))

	fmt.Println("Listening on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
