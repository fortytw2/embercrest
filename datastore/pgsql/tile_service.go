package pgsql

import (
	"github.com/fortytw2/embercrest/game"
	"github.com/jmoiron/sqlx"
)

// TileService is a Tile Service backed by Postgres
type TileService struct {
	db *sqlx.DB
}

// NewTileService returns a new Tile Service from the given database handle
func NewTileService(db *sqlx.DB) *TileService {
	return &TileService{
		db: db,
	}
}

// CreateTile adds a tile to the database
func (tsvc *TileService) CreateTile(c *game.Tile) (err error) {
	return
}

// GetTilees returns all tilees currently in Embercrest
func (tsvc *TileService) GetTiles() (tilees []game.Tile, err error) {
	return
}
