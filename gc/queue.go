package gc

import "github.com/fortytw2/embercrest/user"

// Queue defines the interface used in the Matchmaking queue
type Queue interface {
	Enter(*user.User) error
	Exit(string) error
	Pop() (string, error)
	Length() (int64, error)
}
