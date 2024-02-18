package usergrp

import (
	"ardanlabs/service/business/core/user"
	"ardanlabs/service/business/web/v1/auth"
	"ardanlabs/service/business/web/v1/response"
	"ardanlabs/service/foundation/web"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type handlers struct {
	user *user.Core
	auth *auth.Auth
}

func new(user *user.Core, auth *auth.Auth) *handlers {
	return &handlers{
		user: user,
		auth: auth,
	}
}

// create adds a new user to the system.
func (h *handlers) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}
