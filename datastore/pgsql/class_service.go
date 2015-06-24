package pgsql

import (
	"github.com/fortytw2/embercrest/game"
	"github.com/jmoiron/sqlx"
)

var (
	createClassQuery = `INSERT INTO classes ()`
)

// ClassService is a Class Service backed by Postgres
type ClassService struct {
	db *sqlx.DB
}

// NewClassService returns a new Class Service from the given database handle
func NewClassService(db *sqlx.DB) *ClassService {
	return &ClassService{
		db: db,
	}
}

// CreateClass adds a class to the database
func (cs *ClassService) CreateClass(c *game.Class) (err error) {
	_, err = cs.db.NamedQuery(`INSERT INTO classes
												(name, initialcost, levelcost,
												 hpgrowth, attackgrowth, defensegrowth, hitgrowth, dodgegrowth, critgrowth,
												 minattackrange, maxattackrange,
												 basehp, baseattack, basedefense, basehit, basedodge, basecrit,
												 hpcap, attackcap, defensecap, hitcap, dodgecap, critcap
												) VALUES
												(:name, :initialcost, :levelcost,
												 :hpgrowth, :attackgrowth, :defensegrowth, :hitgrowth, :dodgegrowth, :critgrowth,
												 :minattackrange, :maxattackrange,
												 :basehp, :baseattack, :basedefense, :basehit, :basedodge, :basecrit,
												 :hpcap, :attackcap, :defensecap, :hitcap, :dodgecap, :critcap);`, *c)
	return
}

// GetClass returns a single class by name
func (cs *ClassService) GetClass(name string) (*game.Class, error) {
	var row *sqlx.Row
	row = cs.db.QueryRowx("SELECT * FROM classes WHERE name = $1;", name)

	var class game.Class
	err := row.StructScan(&class)
	if err != nil {
		return nil, err
	}

	return &class, nil
}

// GetClasses returns all classes currently in Embercrest
func (cs *ClassService) GetClasses() (classes []game.Class, err error) {
	var rows *sqlx.Rows
	rows, err = cs.db.Queryx("SELECT * FROM classes;")
	if err != nil {
		return
	}

	for rows.Next() {
		var class game.Class
		err = rows.StructScan(&class)
		if err != nil {
			return
		}

		classes = append(classes, class)
	}
	return
}
