package tests

import (
	"ardanlabs/service/app/services/sales-api/v1/handlers"
	"ardanlabs/service/app/services/sales-api/v1/handlers/usergrp"
	"ardanlabs/service/business/core/user"
	"ardanlabs/service/business/data/dbtest"
	"ardanlabs/service/business/data/order"
	v1 "ardanlabs/service/business/web/v1"
	"ardanlabs/service/business/web/v1/response"
	"ardanlabs/service/foundation/docker"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"runtime/debug"
	"testing"

	"github.com/google/go-cmp/cmp"
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
		usrs, err := api.User.Query(ctx, user.QueryFilter{}, order.By{Field: user.OrderByName, Direction: order.ASC}, 1, 2)
		if err != nil {
			return seedData{}, fmt.Errorf("seeding users: %w", err)
		}

		sd := seedData{
			users: usrs,
		}

		return sd, nil
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
	return func(t *testing.T) {
		table := []struct {
			name    string
			url     string
			resp    any
			expResp any
		}{
			{
				name: "user",
				url:  "/v1/users?page=1&rows=2&orderBy=user_id,DESC",
				resp: &response.PageDocument[usergrp.AppUser]{},
				expResp: &response.PageDocument[usergrp.AppUser]{
					Page:        1,
					RowsPerPage: 2,
					Total:       len(sd.users),
					Items:       toAppUsers(sd.users),
				},
			},
		}

		for _, tt := range table {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+wt.adminToken)
			wt.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Errorf("%s: should receive a status code of 200 for the response: %d", tt.name, w.Code)
				continue
			}

			if err := json.Unmarshal(w.Body.Bytes(), tt.resp); err != nil {
				t.Errorf("should be able to unmashal the reponse : %s", err)
				continue
			}

			diff := cmp.Diff(tt.resp, tt.expResp)
			if diff != "" {
				t.Error("should get the expected response")
				t.Log("GOT")
				t.Logf("%#v", tt.resp)
				t.Log("EXP")
				t.Logf("%#v", tt.expResp)
				continue
			}
		}
	}

}
