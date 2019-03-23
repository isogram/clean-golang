package entity

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//User data
type User struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"fullname"`
	UserName  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRegister data
type UserRegister struct {
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserLogin data
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRefreshToken data
type UserRefreshToken struct {
	RefreshToken string `json:"refresh_token"  validate:"required"`
}

// UserLoginResponse data
type UserLoginResponse struct {
	*User
	AuthToken
}

// AuthToken ...
// This is what is retured to the user
type AuthToken struct {
	TokenType    string `json:"token_type"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthTokenClaim ...
// This is the cliam object which gets parsed from the authorization header
type AuthTokenClaim struct {
	*jwt.StandardClaims
	UID string `json:"uid"`
}

// AuthRefreshToken ...
// This is the refresh token data based on db
type AuthRefreshToken struct {
	ID           int64
	UserID       int64
	RefreshToken string
	Revoked      int
}
