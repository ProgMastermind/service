package web

import (
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

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
