package pgsql

import (
	"testing"

	"github.com/fortytw2/embercrest/datastore"
	"github.com/fortytw2/embercrest/game"
	"github.com/fortytw2/embercrest/user"
)

var ds *datastore.Datastore

func init() {
	var err error
	ds, err = NewDatastore()
	if err != nil {
		panic(err)
	}
}

func getFakeTiles() []game.Tile {
	t := game.Tile{
		ID:           1,
		Name:         "test_tile",
		Resistance:   1,
		DodgeBonus:   5,
		DefenseBonus: 0,
	}

	return []game.Tile{t}
}

func TestMatchService(t *testing.T) {
	m := game.NewMatch("luke", "anakin", getFakeTiles())
	err := ds.CreateMatch(m)
	if err != nil {
		t.Errorf("create match returned error %s", err)
	}

	match, err := ds.GetMatch(m.UUID)
	if err != nil {
		t.Errorf("get match returned error %s", err)
	}
	if match.UUID != m.UUID {
		t.Error("match UUID does not match inserted match UUID")
	}

	matches, err := ds.GetUsersMatches("luke", true)
	if err != nil {
		t.Errorf("create match returned error %s", err)
	}

	if len(matches) != 1 {
		t.Error("more than one match returned...")
	}

	newPlayers := []string{"leia", "mace"}
	matches[0].Players = newPlayers
	err = ds.UpdateMatch(&matches[0])
	if err != nil {
		t.Errorf("update match returned error %s", err)
	}

	match, err = ds.GetMatch(m.UUID)
	if err != nil {
		t.Errorf("get match returned error %s", err)
	}
	if match.Players[0] != newPlayers[0] && match.Players[1] != newPlayers[1] {
		t.Error("match was not updated to have new players, get your shit together dude")
	}
}

func TestUserService(t *testing.T) {
	u, err := user.CreateUser("luke", "luke@jedi.org", "iminlovewithmysister")
	if err != nil {
		t.Errorf("user.CreateUser returned error %s", err)
	}

	err = ds.CreateUser(u)
	if err != nil {
		t.Errorf("create user returned error %s", err)
	}

	u, err = ds.GetUser("luke")
	if err != nil {
		t.Errorf("get user returned error %s", err)
	}

	u.Username = "darth luke bro"
	err = ds.UpdateUser(u)
	if err != nil {
		t.Errorf("update user returned error %s", err)
	}

	u, err = ds.GetUser("darth luke bro")
	if err != nil {
		t.Errorf("get user returned error %s", err)
	}

	return
}

func TestTileService(t *testing.T) {
	err := ds.CreateTile(&game.Tile{Name: "dirt"})
	if err != nil {
		t.Errorf("create tile returned error %s", err)
	}

	tiles, err := ds.GetTiles()
	if err != nil {
		t.Errorf("get all tiles returned error %s", err)
	}
	if len(tiles) != 1 {
		t.Error("GetTiles returns something other than one tile...")
	}

	return
}
