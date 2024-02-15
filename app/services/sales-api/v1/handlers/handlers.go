package handlers

import (
	"ardanlabs/service/app/services/sales-api/v1/checkgrp"
	"ardanlabs/service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "ardanlabs/service/business/web/v1"
	"ardanlabs/service/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface to add all the routes
func (Routes) Add(app *web.App, apiCfg v1.APIMuxConfig) {
	hackgrp.Routes(app, hackgrp.Config{
		Auth: apiCfg.Auth,
	})

	checkgrp.Routes(app, checkgrp.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
	})
}
