package bolt

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/octavore/naga/service"
	"github.com/stretchr/testify/assert"

	"github.com/ketchuphq/ketchup/proto/ketchup/models"
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
		Token:          proto.String("aToken"),
		HashedPassword: proto.String("defghi"),
	}

	assert.NoError(t, app.UpdateUser(a))
	assert.NoError(t, app.UpdateUser(b))
	assert.EqualError(t, app.UpdateUser(&models.User{
		Email:          proto.String("b@example.com"),
		HashedPassword: proto.String("defghi"),
	}), "user already exists")
	assert.EqualError(t, app.UpdateUser(&models.User{
		Uuid:           proto.String("fake"),
		Email:          proto.String("b@example.com"),
		HashedPassword: proto.String("defghi"),
	}), "user already exists")

	assert.NotNil(t, a.GetUuid())
	assert.NotNil(t, b.GetUuid())

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
		assert.NoError(t, err)
		if !tt.valid {
			assert.Nil(t, user)
			continue
		}
		assert.Equal(t, tt.email, user.GetEmail())
		assert.Equal(t, tt.hashedPassword, user.GetHashedPassword())
		sameUser, err := app.GetUser(user.GetUuid())
		assert.NoError(t, err)
		assert.Equal(t, user, sameUser)
	}

	user, err := app.GetUserByToken("aToken")
	assert.NoError(t, err)
	assert.Equal(t, b, user)

	user, err = app.GetUserByToken("notToken")
	assert.NoError(t, err)
	assert.Nil(t, user)

	users, err := app.ListUsers()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []*models.User{a, b}, users)
}
