package db

import (
	"fmt"

	"github.com/James-D-Wood/blog-api/internal/httputils"
	"github.com/James-D-Wood/blog-api/internal/model"
	"github.com/google/uuid"
)

// DefaultUserMap establishes our list of dummy users
//
// | Username  | Password    | Is Admin |
// | --------- | ----------- | -------- |
// | kishiguro | hailsham    | false    |
// | dsedaris  | emeraldIsle | false    |
// | admin     | password    | true     |
var DefaultUserMap = map[string]*model.User{
	"kishiguro": {
		ID:       assignUUID(),
		Username: "kishiguro",
		Name:     "Kazuo Ishiguro",
		IsAdmin:  false,
	},
	"dsedaris": {
		ID:       assignUUID(),
		Username: "dsedaris",
		Name:     "David Sedaris",
		IsAdmin:  false,
	},
	"admin": {
		ID:       assignUUID(),
		Username: "admin",
		Name:     "James Wood",
		IsAdmin:  true,
	},
}

func assignUUID() string {
	i, _ := uuid.NewV7()
	return i.String()
}

// UserService is an abstraction over database actions that can take place on behalf of a user
type UserService interface {
	AuthenticateUser(id, password string) (string, error)
	FetchUser(username string) (*model.User, error)
}

// InMemoryUserService implements UserService to mock user login functionality that would otherwise be handled by an auth service
type InMemoryUserService struct {
	// Users maps username to user details
	Users map[string]*model.User
}

// AuthenticateUser returns a JWTToken representing the user
func (s *InMemoryUserService) AuthenticateUser(username, password string) (string, error) {
	// a password check is not being executed here as this is just simulating what an auth service may return
	user, err := s.FetchUser(username)
	if err != nil {
		return "", fmt.Errorf("error authenticating user: %s", err)
	}

	token, err := httputils.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("error generating JWT for user: %s", err)
	}

	return token, nil
}

// FetchUser returns a user by username
func (s *InMemoryUserService) FetchUser(username string) (*model.User, error) {
	if user, ok := s.Users[username]; ok {
		return user, nil
	}
	return nil, ErrEntityNotFound
}
