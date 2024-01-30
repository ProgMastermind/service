package handlers

import (
	"ardanlabs/service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "ardanlabs/service/business/web/v1"
	"ardanlabs/service/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgrp.Routes(app)
}
