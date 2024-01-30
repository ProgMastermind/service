package v1

import (
	"ardanlabs/service/foundation/logger"
	"ardanlabs/service/foundation/web"
	"os"
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
	Add(app *web.App, cfg APIMuxConfig)
}

// APIMux constructs a http.Handler with all application from routes defined
func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *web.App {
	app := web.NewApp(cfg.Shutdown)

	routeAdder.Add(app, cfg)

	return app
}
