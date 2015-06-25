package user

import (
	"errors"
	"log"
	"os"

	"github.com/fortytw2/abdi"
)

// ensure abdi has the filesystem key for HMAC
// this key is loaded into env from .env by github.com/joho/godotenv
func init() {
	abdi.Key = []byte(os.Getenv("SECRETKEY"))
	// ensure we have a secret key
	if abdi.Key == nil {
		log.Fatalln("FATAL: secret key not present")
	}
}

// User model
type User struct {
	ID int `json:"-"`

	Username     string `json:"username"`
	Email        string `json:"-"`
	PasswordHash string `json:"-"`

	Elo int `json:"elo"`

	Approved  bool `json:"-"`
	Admin     bool `json:"-"`
	Confirmed bool `json:"-"`
}

var (
	errEmailInvalidOrTaken = errors.New("email is invalid or taken")
	errUsernameTaken       = errors.New("username is invalid or taken")
	errLoginFailure        = errors.New("username or password is not valid")
	errInvalidToken        = errors.New("token invalid or expired")
	errAlreadyConfirmed    = errors.New("user already confirmed")
	errInvalidCredentials  = errors.New("invalid basic auth")
	errUserNotFound        = errors.New("unable to find user")
)

// CreateUser creates a new, validated user
func CreateUser(username string, email string, password string) (user *User, err error) {
	var hash string
	hash, err = abdi.Hash(password)
	if err != nil {
		return
	}

	user = &User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Elo:          1000,
		Approved:     false,
		Admin:        false,
		Confirmed:    false,
	}

	return
}

// CheckPassword checks a users password against the password hash and returns
// a bool and any errors
func (u *User) CheckPassword(password string) bool {
	if err := abdi.Check(password, u.PasswordHash); err != nil {
		return false
	}
	return true
}

// GenConfirmationCode creates a confirmationcode using crypto
func (u *User) GenConfirmationCode() (*string, error) {
	return nil, nil
}

// Confirm the user based on the confirmation code passed
func (u *User) Confirm(cc string) error {
	return nil
}
