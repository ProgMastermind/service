package tests

import (
	"ardanlabs/service/app/services/sales-api/v1/handlers"
	"ardanlabs/service/business/data/dbtest"
	v1 "ardanlabs/service/business/web/v1"
	"ardanlabs/service/foundation/docker"
	"context"
	"net/http"
	"os"
	"os/user"
	"runtime/debug"
	"testing"
)

type seedData struct {
	users []user.User
}

// WebTests hold methods for each subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered
var c *docker.Container

type WebTests struct {
	app        http.Handler
	userToken  string
	adminToken string
}

func Test_Web(t *testing.T) {
	t.Parallel()

	test := dbtest.NewTest(t, c, "Test_web")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}

		test.Teardown()
	}()

	api := test.CoreAPIs

	shutdown := make(chan os.Signal, 1)
	tests := WebTests{
		app: v1.APIMux(v1.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.V1.Auth,
			DB:       test.DB,
		}, handlers.Routes{}),
		userToken:  test.TokenV1("user@example.com", "gophers"),
		adminToken: test.TokenV1("admin@example.com", "gophers"),
	}

	seed := func(ctx context.Context, api dbtest.CoreAPIs) (seedData, error) {
		return seedData{}, nil
	}

	t.Log("Seedng data...")

	sd, err := seed(context.Background(), api)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	t.Run("get200", tests.get200(sd))
}

// ---------------------------------------------------------------------------
func (wt *WebTests) get200(sd seedData) func(t *testing.T) {
	return func(t *testing.T) {}
}
