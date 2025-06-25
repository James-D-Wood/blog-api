package model

import (
	"github.com/James-D-Wood/blog-api/internal/httputils"
	"github.com/golang-jwt/jwt/v5"
)

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

	return token.SignedString(httputils.HMACSecret)
}
