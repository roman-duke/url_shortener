package main

import (
	"fmt"
	"net/http"
	"time"

	// "github.com/a-h/templ"
)

// We generally should have a constructor function (idiomatic Go)

// this is cool for generating static data
// as you can see, any argument passed to the templ component is rendered
// at runtime and we have no way to generate dynamic data. But this is quite
// important as we can actually generate the static assets once during build
// time and then just serve those static assets on that route. This would
// prevent computation for each request - basically the whole idea of SSG.

//============================= SSG-ish ====================================//
// func main() {
// 	http.Handle("/", templ.Handler(appShell("Eren Jaeger", time.Now())))
// 	http.Handle("/404", templ.Handler(notFoundComponent()))

// 	fmt.Println("Listening on http://localhost:3000")
// 	http.ListenAndServe(":3000", nil)
// }

// Now we look towards the idea of ssr (which would allow us generate dynamic content)
func NewNowHanlder(now func() time.Time) NowHandler {
	return NowHandler{Now: now}
}

type NowHandler struct {
	Now func() time.Time
}

func (nh NowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timeComponent(nh.Now()).Render(r.Context(), w)
}

// =============================== SSR-ish ==============================//
func main() {
	http.Handle("/", NewNowHanlder(time.Now))

	fmt.Printf("\033[36m\033[1mListening on http://localhost:8080\033[0m\n")

	http.ListenAndServe(":8080", nil)

}
