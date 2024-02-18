package usergrp

import (
	"ardanlabs/service/business/core/user"
	"ardanlabs/service/foundation/validate"
	"fmt"
	"net/mail"
	"time"
)

// AppUser represents information about an individual user.
type AppUser struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Department   string   `json:"department"`
	Enabled      bool     `json:"enabled"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

func toAppUser(usr user.User) AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return AppUser{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []AppUser {
	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

// AppNewUser defines the data needed to add a new user.
type AppNewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

func toCoreNewUser(app AppNewUser) (user.NewUser, error) {
	var roles []user.Role
	if app.Roles != nil {
		roles = make([]user.Role, len(app.Roles))
		for i, roleStr := range app.Roles {
			role, err := user.ParseRole(roleStr)
			if err != nil {
				return user.NewUser{}, fmt.Errorf("parse: %w", err)
			}
			roles[i] = role
		}
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	nu := user.NewUser{
		Name:            app.Name,
		Email:           *addr,
		Roles:           roles,
		Department:      app.Department,
		Password:        app.Password,
		PasswordConfirm: app.PasswordConfirm,
	}

	return nu, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}
