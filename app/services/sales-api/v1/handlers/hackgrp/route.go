package hackgrp

import (
	"ardanlabs/service/business/web/v1/auth"
	"ardanlabs/service/business/web/v1/mid"
	"ardanlabs/service/foundation/web"
	"net/http"
)

// Config contains all the mandatory systems required by handlers
type Config struct {
	Auth *auth.Auth
}

// Routes add specific routes for this group
func Routes(app *web.App, cfg Config) {
	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)
	app.Handle(http.MethodGet, "/hack", Hack)
	app.Handle(http.MethodGet, "/hackauth", Hack, authen, ruleAdmin)
}
