package pgsql

import (
	"os"

	"github.com/fortytw2/embercrest/datastore"
	"github.com/jmoiron/sqlx"
)

// NewDatastore creates a new datastore backed by PGSQL
func NewDatastore() (ds *datastore.Datastore, err error) {
	var db *sqlx.DB
	db, err = sqlx.Connect("postgres", os.Getenv("DATABASE"))
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	ds = &datastore.Datastore{
		UserService:  NewUserService(db),
		MatchService: NewMatchService(db),
		TileService:  NewTileService(db),
		ClassService: NewClassService(db),
	}

	return
}
