package pgsql

import (
	"encoding/json"
	"errors"

	pg_types "github.com/fortytw2/embercrest/datastore/pgsql/types"
	"github.com/fortytw2/embercrest/game"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
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

// CreateMatch adds a new match to the database
func (ms *MatchService) CreateMatch(match *game.Match) error {
	matchJSON, err := json.Marshal(match)
	if err != nil {
		return err
	}

	_, err = ms.db.Queryx("INSERT INTO matches (uuid, usernames, active, match) VALUES ($1, $2, $3, $4)", match.UUID, pg_types.PGStringSlice(match.Players), true, types.JsonText(matchJSON))
	return err
}

// GetMatch returns a single match by UUID
func (ms *MatchService) GetMatch(uuid string) (match *game.Match, err error) {
	row := ms.db.QueryRowx("SELECT match FROM matches WHERE uuid = $1;", uuid)

	var text types.JsonText
	err = row.Scan(&text)
	if err != nil {
		return
	}

	err = text.Unmarshal(&match)
	if err != nil {
		return
	}

	return
}

// UpdateMatch updates a match in the database from it's UUID
func (ms *MatchService) UpdateMatch(m *game.Match) (err error) {
	var match *game.Match
	match, err = ms.GetMatch(m.UUID)
	if err != nil {
		return
	}
	if match.UUID != m.UUID {
		err = errors.New("something went disasterously wrong")
		return
	}

	matchJSON, err := json.Marshal(m)
	if err != nil {
		return err
	}

	isActive := false
	if m.State != game.Ended {
		isActive = true
	}

	_, err = ms.db.Queryx("UPDATE matches SET usernames = $2, active = $3, match = $4 WHERE uuid = $1", m.UUID, pg_types.PGStringSlice(m.Players), isActive, types.JsonText(matchJSON))

	return
}

// GetUsersMatches returns all of a users matches, active or inactive
func (ms *MatchService) GetUsersMatches(username string, active bool) (matches []game.Match, err error) {
	var rows *sqlx.Rows
	rows, err = ms.db.Queryx("SELECT match FROM matches WHERE $1 = ANY(usernames) AND active = $2;", username, active)
	if err != nil {
		return
	}

	for rows.Next() {
		var match game.Match
		var text types.JsonText

		err = rows.Scan(&text)
		if err != nil {
			return
		}

		err = text.Unmarshal(&match)
		if err != nil {
			return
		}
		matches = append(matches, match)
	}

	return
}

// ActiveMatches returns all active matches
func (ms *MatchService) ActiveMatches() (matches []game.Match, err error) {
	var rows *sqlx.Rows
	rows, err = ms.db.Queryx("SELECT match FROM matches WHERE active = true;")
	if err != nil {
		return
	}

	for rows.Next() {
		var match game.Match
		var text types.JsonText

		err = rows.Scan(&text)
		if err != nil {
			return
		}

		err = text.Unmarshal(&match)
		if err != nil {
			return
		}
		matches = append(matches, match)
	}

	return
}
