package repository

import (
	"E-Commerce/api-gateway/internal/auth"
	"errors"
)

type AuthRepository interface {
	CreateUser(user *auth.User) error
	GetUserByEmail(email string) (*auth.User, error)
}

type authRepository struct {
	users map[string]*auth.User
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		users: make(map[string]*auth.User),
	}
}

func (r *authRepository) CreateUser(user *auth.User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *authRepository) GetUserByEmail(email string) (*auth.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
