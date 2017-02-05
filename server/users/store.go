package users

import (
	"github.com/octavore/ketchup/db"
	"github.com/octavore/ketchup/proto/ketchup/models"
)

type userStore struct{ db.Backend }

// Create a user from email and hashedPassword.
// Wraps DB for nagax.users module.
func (s *userStore) Create(email, hashedPassword string) (string, error) {
	user := &models.User{
		Email:          &email,
		HashedPassword: &hashedPassword,
	}
	// todo: whitelist?
	err := s.UpdateUser(user)
	if err != nil {
		return "", err
	}
	return user.GetUuid(), nil
}

// Get a user id based on email.
// Wraps DB for nagax.users module.
func (s *userStore) Get(email string) (id, hashedPassword string, err error) {
	u, err := s.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}
	return u.GetUuid(), u.GetHashedPassword(), nil
}
