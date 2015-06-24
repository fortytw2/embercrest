package pgsql

import (
	"testing"

	"github.com/fortytw2/embercrest/datastore"
	"github.com/fortytw2/embercrest/game"
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

func TestClassService(t *testing.T) {
	err := ds.CreateClass(&game.Class{Name: "jedi"})
	if err != nil {
		t.Errorf("create class returned error %s", err)
	}

	class, err := ds.GetClass("jedi")
	if err != nil {
		t.Errorf("get class returned error %s", err)
	}
	if class.Name != "jedi" {
		t.Errorf("get class returned class with incorrect name %s", class.Name)
	}

	classes, err := ds.GetClasses()
	if err != nil {
		t.Errorf("get all classes returned error %s", err)
	}
	if len(classes) != 1 {
		t.Error("getClasses returns something other than one class? really?")
	}

	return
}
