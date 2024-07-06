package yume_test

import (
	"context"
	"testing"

	"encore.app/yume"
)

func TestTestYume(t *testing.T) {
	t.Parallel()
	return
}

func TestNewUser(t *testing.T) {
	t.Parallel()
	user, err := yume.NewUser(context.Background(), &yume.NewUserParams{ID: "test", USERNAME: "test"})
	if err != nil {
		t.Errorf("unexpected error from NewUser: %v", err)
	}
	if user == nil {
		t.Errorf("unexpected nil user")
	}
	// check if user exists in db
	res, err := yume.CheckUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf("unexpected error from CheckUser: %v", err)
	}
	if res == nil {
		t.Errorf("unexpected nil result from CheckUser")
	}
	if !res.Exists {
		t.Errorf("unexpected user not exist error")
	}
}
