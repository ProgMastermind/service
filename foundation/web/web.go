package web

import (
	"context"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

// A Handler is a type that handles a http request within our own little mini framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is hte entrypoint into our application and what configures our context
// object for each of out http handlers. Feel free to add any configuration
// data/logic on this app struct

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux

func (a *App) Handle(method string, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {

		// ADD ANY LOGIC HERE

		if err := handler(r.Context(), w, r); err != nil {
			// NOT SURE WHAT TO DO
			return
		}

		// ADD ANY LOGIC HERE

	}

	a.ContextMux.Handle(method, path, h)
}
