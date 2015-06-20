package game

// Tile is an object on a coordinate in the map, making up the terrain
type Tile struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Resistance   int    `json:"move_resistance"`
	DodgeBonus   int    `json:"dodge_bonus"`
	DefenseBonus int    `json:"def_bonus"`
}

// Point represents a cartesian coordinate
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
