package pgsql

import (
	"github.com/fortytw2/embercrest/datastore"
	"github.com/fortytw2/embercrest/game"
	"github.com/jmoiron/sqlx"
)

// MatchService is a Match Service backed by Postgres
type MatchService struct {
	db *sqlx.DB
}

// NewMatchService returns a new Match Service from the given database handle
func NewMatchService(db *sqlx.DB) *MatchService {
	return &MatchService{
		db: db,
	}
}

func (msvc *MatchService) CreateMatch(match *game.Match) (err error) {
	return
}

func (msvc *MatchService) GetMatch(uuid string) (match *game.Match, err error) {
	return
}

func (msvc *MatchService) UpdateMatch(m *game.Match) (err error) {
	return
}

func (msvc *MatchService) SearchMatches(ms *datastore.MatchSearchParams) (matches []game.Match, err error) {
	return
}

func (msvc *MatchService) ActiveMatches() (matches []game.Match, err error) {
	return
}
