package models

import "github.com/golang-jwt/jwt/v5"

type Credentials struct {
	Email string `json:"email" validate:"required,email"`
	Pass  string `json:"password" validate:"required"`
}

type Claims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
