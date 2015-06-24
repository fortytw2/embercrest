package game

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/satori/go.uuid"
)

type MatchState string

const (
	InDrafting MatchState = "drafting in progress"
	InPlay     MatchState = "game in progress"
	Ended      MatchState = "game ended"
)

// Match contains all data that occurs during a game between two players
type Match struct {
	UUID      string    `json:"id"`
	StartTime time.Time `json:"start"`
	LastTurn  time.Time `json:"-"`
	EndTime   time.Time `json:"end,omitempty"`
	Players   []string  `json:"players"`

	Map Map `json:"map"`

	Funds        map[string]int      `json:"funds"`
	Armies       map[string][]Unit   `json:"armies,omitempty"`
	DraftHistory map[string][]Unit   `json:"draft,omitempty"`
	ArmyHistory  map[string][][]Unit `json:"army_history,omitempty"`

	State      MatchState `json:"match_state"`
	PlayerTurn string     `json:"current_player_turn"` // ID of the user whose turn it is
	TurnNumber int        `json:"total_turns"`         // number of turns so far
	Victor     string     `json:"victor,omitempty"`
}

// NewMatch creates a new Match between two users
func NewMatch(userID1 string, userID2 string, tiles []Tile) *Match {
	return &Match{
		UUID:         uuid.NewV4().String(),
		StartTime:    time.Now(),
		Funds:        map[string]int{userID1: 2000, userID2: 2000},
		Players:      []string{userID1, userID2},
		Map:          GenerateMap(tiles),
		Armies:       make(map[string][]Unit, 2),
		DraftHistory: make(map[string][]Unit, 2),
		ArmyHistory:  make(map[string][][]Unit, 2),
		State:        InDrafting,
		PlayerTurn:   userID1,
	}

}

// Draft handles picking of units in a back and forth manner
func (m *Match) Draft(class string, level int, c *Class) error {

	var unitLoc Point
	if m.PlayerTurn == m.Players[0] {
		unitLoc.X = 20 - len(m.DraftHistory[m.PlayerTurn])
		unitLoc.Y = 20
	} else {
		unitLoc.X = len(m.DraftHistory[m.PlayerTurn])
		unitLoc.Y = 0
	}

	unit := CreateUnit(c, level, unitLoc)
	if m.Funds[m.PlayerTurn] < unit.TotalCost {
		return errors.New("not enough funds for unit")
	}
	m.Funds[m.PlayerTurn] = m.Funds[m.PlayerTurn] - unit.TotalCost
	m.DraftHistory[m.PlayerTurn] = append(m.DraftHistory[m.PlayerTurn], *unit)
	m.Armies[m.PlayerTurn] = append(m.Armies[m.PlayerTurn], *unit)

	m.NextTurn()

	return nil
}

// Move represents the movement of one unit
type Move struct {
	UnitUUID string  `json:"unit_id"`
	Path     []Point `json:"path"`
	Attack   string  `json:"attack,omitempty"`
}

// MakeMoves executes the moves on the match
func (m *Match) MakeMoves(moves []Move) error {
	for _, move := range moves {
		err := m.validateMove(move)
		if err != nil {
			return err
		}
	}

	for _, move := range moves {
		m.makeMove(move)
	}

	// advance to the next turn, changing state as needed
	m.NextTurn()

	return nil
}

func (m *Match) makeMove(move Move) error {
	return nil
}

// NextTurn advances the match to the next players turn
func (m *Match) NextTurn() {
	if m.testDraftOver() {
		m.State = InPlay
	}
	if m.testGameOver() {
		m.State = Ended
	}

	if m.PlayerTurn == m.Players[0] {
		m.PlayerTurn = m.Players[1]
	} else {
		m.PlayerTurn = m.Players[0]
	}
	m.TurnNumber++
	m.LastTurn = time.Now()
}

func (m *Match) validateMove(move Move) error {
	unit, err := m.findUnit(move.UnitUUID)
	if err != nil {
		return err
	}
	err = m.checkPath(unit, move.Path)
	if err != nil {
		return err
	}
	err = m.checkAttack(unit, move.Attack)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) findUnit(id string) (*Unit, error) {
	var u Unit

	for _, unit := range m.Armies[m.PlayerTurn] {
		if unit.UUID == id {
			u = unit
		}
	}
	if &u != nil {
		return nil, fmt.Errorf("unable to find unit, %s", id)
	}

	return &u, nil
}

func (m *Match) checkAttack(unit *Unit, attackUnit string) error {
	u, err := m.findUnit(attackUnit)
	if err != nil {
		return err
	}
	// ensure the unit being attacked is between the min and max attack range
	if unit.Coordinate.X+unit.MinAttackRange <= u.Coordinate.X && unit.Coordinate.X+unit.MaxAttackRange >= u.Coordinate.X {

	} else {
		return fmt.Errorf("unit %s not within range to attack unit %s", unit.UUID, attackUnit)
	}

	return nil
}

func (m *Match) checkPath(unit *Unit, path []Point) error {
	stepsLeft := unit.Move
	// start moving at the unit's last position
	prevPoint := unit.Coordinate

	for _, point := range path {
		if stepsLeft <= 0 {
			return fmt.Errorf("path for unit %s too long, check terrain resistance", unit.UUID)
		}

		if abs(prevPoint.X-point.X) <= 1 && abs(prevPoint.Y-point.Y) <= 1 {
			stepsLeft -= m.Map[point.X][point.Y].Resistance
		} else {
			return fmt.Errorf("path is not correct, point %v is not next to point %v", prevPoint, point)
		}
	}
	return nil

}

func (m *Match) testDraftOver() bool {
	return (m.Funds[m.Players[0]] == 0 && m.Funds[m.Players[1]] == 0)
}

func (m *Match) testGameOver() bool {
	if (len(m.Armies[m.Players[0]]) == 0 || len(m.Armies[m.Players[1]]) == 0) && m.State == InPlay {
		log.Println("mm: match", m.UUID, "ended")
		return true
	}
	return false
}

func abs(a int) int {
	if a < 0 {
		a *= -1
	}
	return a
}
