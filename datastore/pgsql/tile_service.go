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
func (ts *TileService) CreateTile(t *game.Tile) (err error) {
	_, err = ts.db.NamedQuery(`INSERT INTO tiles
												(name, resistance, defensebonus, dodgebonus) VALUES
												(:name, :resistance, :defensebonus, :dodgebonus);`, *t)
	return
}

// GetTiles returns all tiles
func (ts *TileService) GetTiles() (tiles []game.Tile, err error) {
	var rows *sqlx.Rows
	rows, err = ts.db.Queryx("SELECT * FROM tiles;")
	if err != nil {
		return
	}

	for rows.Next() {
		var tile game.Tile
		err = rows.StructScan(&tile)
		if err != nil {
			return
		}

		tiles = append(tiles, tile)
	}
	return
}
