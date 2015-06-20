package game

type Class struct {
	Name string `json:"class_name"`

	// Class costs
	InitialCost int `json:"initial_cost"`
	LevelCost   int `json:"level_cost"`

	// Class Growth Percentages
	HPGrowth      int `json:"hp_growth"`
	AttackGrowth  int `json:"attack_growth"`
	DefenseGrowth int `json:"defense_growth"`
	HitGrowth     int `json:"hit_growth"`
	DodgeGrowth   int `json:"dodge_growth"`
	CritGrowth    int `json:"crit_growth"`

	// Static Class Stats
	MinAttackRange int `json:"min_attack_range"`
	MaxAttackRange int `json:"max_attack_range"`
	Move           int `json:"move"`

	// Base Class Stats
	BaseHP      int `json:"base_hp"`
	BaseAttack  int `json:"base_attack"`
	BaseDefense int `json:"base_defense"`
	BaseHit     int `json:"base_hit"`
	BaseDodge   int `json:"base_dodge"`
	BaseCrit    int `json:"base_crit"`

	HPCap      int `json:"hp_cap"`
	AttackCap  int `json:"attack_cap"`
	DefenseCap int `json:"defense_cap"`
	HitCap     int `json:"hit_cap"`
	DodgeCap   int `json:"dodge_cap"`
	CritCap    int `json:"crit_cap"`
}
