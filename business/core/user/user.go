package user

import (
	"ardanlabs/service/business/data/order"
	"ardanlabs/service/foundation/logger"
	"context"
	"net/mail"

	"github.com/google/uuid"
)

// Storer interface declares the behavior this package needs to perists and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// ===============================================================================

// Core manages the set of APIs for user access.
type Core struct {
	log    *logger.Logger
	storer Storer
}

// NewCore constructs a user core API for use.
func NewCore(log *logger.Logger, storer Storer) *Core {
	return &Core{
		log:    log,
		storer: storer,
	}
}

// Create adds a new user to the system.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {

	return User{}, nil
}
