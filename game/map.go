package game

import (
	"encoding/json"
	"math/rand"
	"time"
)

// Map is a 2d grid of tiles that implements json.Marshaler
type Map [][]Tile

// GenerateMap creates a ignorantly random map according to tile distributions
func GenerateMap(tiles []Tile) Map {

	// probabilities of getting each TerrainElement
	water := 0.1
	forest := 0.3
	mountain := 0.5

	//randomly generate size of the map
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	mapSize := 20

	// create the array
	lm := make([][]Tile, mapSize)

	// make it 2D
	for i := range lm {
		lm[i] = make([]Tile, mapSize)
	}

	// fill the array with TerrainElements
	for i, row := range lm {
		for j := range row {
			TerrainElementType := r.Float64()
			if TerrainElementType <= water {
				lm[i][j] = tiles[0]
			} else if TerrainElementType <= forest {
				lm[i][j] = tiles[1]
			} else if TerrainElementType <= mountain {
				lm[i][j] = tiles[2]
			} else {
				lm[i][j] = tiles[3]
			}
		}
	}

	return lm
}

// MarshalJSON makes map a json.Marshaler and ensures that the json encoded Match
// is simply a [][]int instead of a huge mess
func (m Map) MarshalJSON() ([]byte, error) {
	// create the array
	lm := make([][]int, 20)

	// make it 2D
	for i := range lm {
		lm[i] = make([]int, 20)
	}

	for i, row := range m {
		for j := range row {
			lm[i][j] = m[i][j].ID
		}
	}
	js, err := json.Marshal(lm)
	if err != nil {
		return nil, err
	}

	return js, nil
}
