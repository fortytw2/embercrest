package game

import (
	"math/rand"

	"github.com/satori/go.uuid"
)

// A Unit is an instance of a class
type Unit struct {
	ClassName      string `json:"class"`
	TotalCost      int    `json:"cost"`
	UUID           string `json:"id"`
	HP             int    `json:"hp"`
	Attack         int    `json:"attack"`
	Defense        int    `json:"defense"`
	Hit            int    `json:"hit_chance"`
	Dodge          int    `json:"dodge_chance"`
	Crit           int    `json:"crit_chance"`
	MinAttackRange int    `json:"min_attack_range"`
	MaxAttackRange int    `json:"max_attack_range"`
	Move           int    `json:"move"`
	Coordinate     Point  `json:"coordinate"`
}

// CreateUnit creates the given unit by the class stats and level
func CreateUnit(class *Class, level int, loc Point) *Unit {
	return &Unit{
		ClassName:      class.Name,
		TotalCost:      class.InitialCost + class.LevelCost*level,
		UUID:           uuid.NewV4().String(),
		HP:             rollStat(class.BaseHP, class.HPGrowth, class.HPCap, level),
		Attack:         rollStat(class.BaseAttack, class.AttackGrowth, class.AttackCap, level),
		Defense:        rollStat(class.BaseDefense, class.DefenseGrowth, class.DefenseCap, level),
		Hit:            rollStat(class.BaseHit, class.HitGrowth, class.HitCap, level),
		Dodge:          rollStat(class.BaseDodge, class.DodgeGrowth, class.DodgeCap, level),
		Crit:           rollStat(class.BaseCrit, class.CritGrowth, class.CritCap, level),
		MinAttackRange: class.MinAttackRange,
		MaxAttackRange: class.MaxAttackRange,
		Move:           class.Move,
		Coordinate:     loc,
	}
}

// roll a stat
func rollStat(base int, chance int, statCap int, level int) int {
	var added int
	for i := 0; int(i) < level; i++ {
		r := rand.Intn(100)
		if chance > r && (base+added <= statCap) {
			added++
		}
	}

	return base + added
}
