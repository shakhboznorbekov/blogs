package auth

import "github.com/golang-jwt/jwt"

type SignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type GenerateToken struct {
	Username string
	Role     string
}

type TokenData struct {
	Username string
	UserId   string
}
type AuthResponse struct {
	Token string `json:"token"`
}
