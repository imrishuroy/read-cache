package auth

import (
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthService provides authentication services
type AuthService struct {
	// DB *gorm.DB
	FireAuth *auth.Client
}

// Login authenticates a user with the provided credentials and returns a Firebase custom token
func (s *AuthService) Login(email, password string) (string, error) {

	// Get the user from the database

	user := User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: password,
	}

	//s.FireAuth.CustomTokenWithClaims()

	authUser, err := s.FireAuth.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		return "", err
	}

	claims := map[string]interface{}{
		"user_email": authUser.Email,
	}

	token, err := s.FireAuth.CustomTokenWithClaims(context.Background(), authUser.UID, claims)
	if err != nil {
		return "", nil
	}

	return token, nil
}

func (s *AuthService) Register(email, password string) (string, error) {

	//uid := uuid.New().String()

	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		//PhoneNumber()
		Password(password).
		//DisplayName()
		//PhotoURL()
		Disabled(false)

	user, err := s.FireAuth.CreateUser(context.Background(), params)
	if err != nil {
		return "", err
	}

	claims := map[string]interface{}{
		"role": "user",
	}

	customToken, err := s.FireAuth.CustomTokenWithClaims(context.Background(), user.UID, claims)
	if err != nil {
		return "", err
	}
	return customToken, nil

}
