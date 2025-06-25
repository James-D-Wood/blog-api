package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// this is insecure - ideally this secret should be managed by another secret management service
var HMACSecret = []byte("7e59e2c4-51a1-11f0-a636-de64e30f34bb")

type User struct {
	ID       string
	Username string
	Name     string
	IsAdmin  bool
}

// this is insecure - ideally we should be expiring our JWTs
func (user *User) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.ID,
			"is_admin": user.IsAdmin,
		},
	)

	return token.SignedString(HMACSecret)
}
