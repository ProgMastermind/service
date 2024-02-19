package usergrp

import (
	"ardanlabs/service/business/core/user"
	"ardanlabs/service/business/core/user/stores/userdb"
	"ardanlabs/service/business/web/v1/auth"
	"ardanlabs/service/business/web/v1/mid"
	"ardanlabs/service/foundation/logger"
	"ardanlabs/service/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

// Routes adds specific routes for this group
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	usrCore := user.NewCore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))

	hdl := new(usrCore, cfg.Auth)
	app.Handle(http.MethodPost, version, "/users", hdl.create)
	app.Handle(http.MethodPost, version, "/usersauth", hdl.create, authen, ruleAdmin)

}
