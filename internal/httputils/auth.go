package httputils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// this is insecure - ideally this secret should be managed by another secret management service
var HMACSecret = []byte("7e59e2c4-51a1-11f0-a636-de64e30f34bb")

type AuthClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
}

func DecodeBasicAuth(r *http.Request) (username, password string, err error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", "", errors.New("no Authorization header provided")
	}

	if len(header) < 6 || !strings.EqualFold(header[:6], "basic ") {
		return "", "", errors.New("basic auth not detected")
	}

	b64String := header[6:]
	b, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return "", "", fmt.Errorf("problem decoding auth header: %s", err)
	}

	components := strings.Split(string(b), ":")
	if len(components) != 2 {
		return "", "", fmt.Errorf("found %d components in basic auth header - expected 2", len(components))
	}

	return components[0], components[1], nil
}

func DecodeBearerAuth(r *http.Request) (token string, err error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("no Authorization header provided")
	}

	if len(header) < 7 || !strings.EqualFold(header[:7], "bearer ") {
		return "", errors.New("bearer auth not detected")
	}

	token = header[7:]
	return token, nil
}

// ExtractJWTClaims verifies that the token was signed by this server and extracts claims about user ID
func ExtractJWTClaims(token string, claims any) error {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return HMACSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return fmt.Errorf("could not parse JWT token: %s", err)
	}

	if jwtClaims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		b, err := json.Marshal(jwtClaims)
		if err != nil {
			return fmt.Errorf("could not parse JWT claims: %s", err)
		}
		err = json.Unmarshal(b, claims)
		if err != nil {
			return fmt.Errorf("could not parse JWT claims: %s", err)
		}
	} else {
		return errors.New("could not parse JWT claims")
	}

	return nil
}
