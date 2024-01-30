package handlers

import (
	"ardanlabs/service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "ardanlabs/service/business/web/v1"

	"github.com/dimfeld/httptreemux/v5"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(mux *httptreemux.ContextMux, cfg v1.APIMuxConfig) {
	hackgrp.Routes(mux)
}
