package user_test

import (
	"ardanlabs/service/business/core/user"
	"ardanlabs/service/business/data/dbtest"
	"ardanlabs/service/foundation/docker"
	"context"
	"fmt"
	"net/mail"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/google/uuid"
)

var c *docker.Container

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error

	c, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}
	defer dbtest.StopDB(c)

	return m.Run(), nil
}

func Test_User(t *testing.T) {
	t.Run("crud", crud)
}

func crud(t *testing.T) {

	seed := func(ctx context.Context, usrCore *user.Core) ([]user.User, error) {
		return []user.User{}, nil
	}

	test := dbtest.NewTest(t, c, "Test_User/crud")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.Teardown()
	}()

	api := test.CoreAPIs

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Log("Go seeding ...")

	_, err := seed(ctx, api.User)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	email, err := mail.ParseAddress("jacob@ardanlabs.com")
	if err != nil {
		t.Fatalf("Should be able to parse email: %s.", err)
	}

	nu := user.NewUser{
		Name:            "Jack Smith",
		Email:           *email,
		Roles:           []user.Role{user.RoleAdmin},
		Department:      "MARK",
		Password:        "123",
		PasswordConfirm: "123",
	}

	usr, err := api.User.Create(context.Background(), nu)
	if err != nil {
		t.Fatalf("Should be able to create user : %s. ", err)
	}

	if usr.ID == uuid.Nil {
		t.Error("Should have a valid ID.")
	}

	if usr.Name != nu.Name {
		t.Error("Should have the correct name.")
		t.Errorf("GOT: %s", usr.Name)
		t.Errorf("EXP: %s", nu.Name)
	}
}
