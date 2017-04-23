package bolt

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"

	"github.com/octavore/ketchup/proto/ketchup/models"
)

func TestUser(t *testing.T) {
	app := &Module{}
	stop := service.New(app).StartForTest()
	defer stop()
	a := &models.User{
		Email:          proto.String("a@example.com"),
		HashedPassword: proto.String("abcdef"),
	}
	b := &models.User{
		Email:          proto.String("b@example.com"),
		HashedPassword: proto.String("defghi"),
	}
	err := app.UpdateUser(a)
	if err != nil {
		t.Fatal(err)
	}
	err = app.UpdateUser(b)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		email          string
		hashedPassword string
		valid          bool
	}{
		{"a@example.com", "abcdef", true},
		{"b@example.com", "defghi", true},
		{"c@example.com", "defghi", false},
	}
	for _, tt := range tests {
		user, err := app.GetUserByEmail(tt.email)
		if err != nil {
			t.Fatal(err)
		}
		if !tt.valid && user != nil {
			t.Error("unexpected user:", user)
		}
		if tt.valid && user.GetEmail() != tt.email {
			t.Errorf("unexpected email %s; wanted %s", user.GetEmail(), tt.email)
		}
		if tt.valid && user.GetHashedPassword() != tt.hashedPassword {
			t.Errorf("unexpected pwd %s; wanted %s", user.GetHashedPassword(), tt.hashedPassword)
		}
	}
}
