package auth

import (
    "github.com/jmoiron/sqlx"
)

type AuthRepository interface {
    CreateUser(user *User) error
    GetUserByEmail(email string) (*User, error)
}

type authRepository struct {
    db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
    return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *User) error {
    query := `INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, user.ID, user.Email, user.Password, user.Role)
    return err
}

func (r *authRepository) GetUserByEmail(email string) (*User, error) {
    var user User
    err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
    if err != nil {
        return nil, err
    }
    return &user, nil
}