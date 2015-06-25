package pgsql

import (
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

// CreateUser adds a user to the datastore
func (us *UserService) CreateUser(u *user.User) (err error) {
	_, err = us.db.NamedQuery("INSERT INTO users (username, email, passwordhash, elo, approved, admin, confirmed) VALUES (:username, :email, :passwordhash, :elo, :approved, :admin, :confirmed)", u)
	return err
}

// UpdateUser updates a database user by ID
func (us *UserService) UpdateUser(u *user.User) (err error) {
	_, err = us.db.NamedQuery(`UPDATE users SET username = :username,
																					 passwordhash = :passwordhash,
																				 	 email = :email,
																					 elo = :elo,
																					 approved = :approved,
																					 admin = :admin,
																					 confirmed = :confirmed
																					 WHERE id = :id`, u)

	return
}

// GetUser returns a user by their username
func (us *UserService) GetUser(username string) (u *user.User, err error) {
	row := us.db.QueryRowx("SELECT * FROM users WHERE username = $1;", username)

	var scanUser user.User
	err = row.StructScan(&scanUser)
	if err != nil {
		return
	}
	u = &scanUser

	return
}

// GetLeaderboard returns the top 25 users by Elo
func (us *UserService) GetLeaderboard() (users []user.User, err error) {
	var rows *sqlx.Rows
	rows, err = us.db.Queryx("SELECT * FROM users SORT BY elo DESC LIMIT 25;")
	if err != nil {
		return
	}

	for rows.Next() {
		var u user.User
		rows.StructScan(&u)
		users = append(users, u)
	}

	return
}
