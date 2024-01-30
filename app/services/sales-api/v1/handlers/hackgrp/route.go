package hackgrp

import (
	"ardanlabs/service/foundation/web"
	"net/http"
)

// Routes add specific routes for this group
func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
