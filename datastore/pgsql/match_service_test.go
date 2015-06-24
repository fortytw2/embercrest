package pgsql

import (
	"testing"

	"github.com/fortytw2/embercrest/game"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ms *MatchService
var uuid string

func init() {
	db, err := sqlx.Connect("postgres", "user=postgres password=root dbname=embercrest sslmode=disable")
	if err != nil {
		panic(err)
	}
	ms = NewMatchService(db)
}

func getTiles() []game.Tile {
	t := game.Tile{
		ID:           1,
		Name:         "test_tile",
		Resistance:   1,
		DodgeBonus:   5,
		DefenseBonus: 0,
	}

	return []game.Tile{t}
}

func TestCreateMatch(t *testing.T) {
	m := game.NewMatch("luke", "anakin", getTiles())
	err := ms.CreateMatch(m)
	if err != nil {
		t.Errorf("create match returned error %s", err)
	}

	match, err := ms.GetMatch(m.UUID)
	if err != nil {
		t.Errorf("get match returned error %s", err)
	}
	if match.UUID != m.UUID {
		t.Error("match UUID does not match inserted match UUID")
	}

	matches, err := ms.GetUsersMatches("luke", true)
	if err != nil {
		t.Errorf("create match returned error %s", err)
	}

	if len(matches) != 1 {
		t.Error("more than one match returned...")
	}

	newPlayers := []string{"leia", "mace"}
	matches[0].Players = newPlayers
	err = ms.UpdateMatch(&matches[0])
	if err != nil {
		t.Errorf("update match returned error %s", err)
	}

	match, err = ms.GetMatch(m.UUID)
	if err != nil {
		t.Errorf("get match returned error %s", err)
	}
	if match.Players[0] != newPlayers[0] && match.Players[1] != newPlayers[1] {
		t.Error("match was not updated to have new players, get your shit together dude")

	}

}
