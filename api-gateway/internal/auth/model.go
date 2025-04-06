package auth

type User struct {
    ID       string `json:"id" db:"id"`
    Email    string `json:"email" db:"email"`
    Password string `json:"-" db:"password_hash"`
    Role     string `json:"role" db:"role"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}