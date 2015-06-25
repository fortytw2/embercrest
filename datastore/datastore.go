package datastore

import (
	"github.com/fortytw2/embercrest/game"
	"github.com/fortytw2/embercrest/user"
)

// A Datastore holds all storage services
type Datastore struct {
	UserService
	MatchService
	ClassService
	TileService
}

// UserService provides a wrapper around user persistance functions
type UserService interface {
	CreateUser(*user.User) error
	UpdateUser(*user.User) error
	GetUser(username string) (*user.User, error)

	GetLeaderboard() ([]user.User, error)
}

// MatchService provides a wrapper around match persistance functions
type MatchService interface {
	CreateMatch(*game.Match) error
	GetMatch(string) (*game.Match, error)
	UpdateMatch(*game.Match) error

	GetUsersMatches(username string, active bool) ([]game.Match, error)
	ActiveMatches() ([]game.Match, error)
}

// ClassService provides a wrapper around class persistance functions
type ClassService interface {
	CreateClass(*game.Class) error
	GetClass(string) (*game.Class, error)
	GetClasses() ([]game.Class, error)
}

// TileService provides a wrapper around tile persistance functions
type TileService interface {
	CreateTile(*game.Tile) error
	GetTiles() ([]game.Tile, error)
}
