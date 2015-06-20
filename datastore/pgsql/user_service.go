package pgsql

import (
	"github.com/fortytw2/embercrest/game"
	"github.com/fortytw2/embercrest/user"
	"github.com/jmoiron/sqlx"
)

// UserService is a User Service backed by Postgres
type UserService struct {
	db *sqlx.DB
}

// NewUserService returns a new User Service from the given database handle
func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (usvc *UserService) CreateUser(u *user.User) (err error) {
	return
}
func (usvc *UserService) UpdateUser(u *user.User) (err error) {
	return
}
func (usvc *UserService) GetUser(username string) (u *user.User, err error) {
	return
}

func (usvc *UserService) GetUserMatches(username string) (matches []game.Match, err error) {
	return
}

func (usvc *UserService) GetLeaderboard() (users []user.User, err error) {
	return
}
