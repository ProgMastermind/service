package v1

import (
	"ardanlabs/service/foundation/logger"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

// APIMuxConfig contains all the mandatory systems required by handlers
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

// RouteAdder defines behaviour that sets the routes to bind for an instance
// of the service

type RouteAdder interface {
	Add(mux *httptreemux.ContextMux, cfg APIMuxConfig)
}

// APIMux constructs a http.Handler with all application from routes defined
func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	routeAdder.Add(mux, cfg)

	return mux
}
