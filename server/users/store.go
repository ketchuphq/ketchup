package users

import (
	"github.com/octavore/ketchup/proto/ketchup/models"
)

// Create a user from email and hashedPassword. For nagax.users module.
func (m *Module) Create(email, hashedPassword string) (string, error) {
	user := &models.User{
		Email:          &email,
		HashedPassword: &hashedPassword,
	}
	// todo: whitelist?
	err := m.DB.UpdateUser(user)
	if err != nil {
		return "", err
	}
	return user.GetUuid(), nil
}

// Get a user id based on email. For the nagax.users module.
func (m *Module) Get(email string) (id, hashedPassword string, err error) {
	u, err := m.DB.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}
	return u.GetUuid(), u.GetHashedPassword(), nil
}
