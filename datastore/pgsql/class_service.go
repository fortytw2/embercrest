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
func (csvc *ClassService) CreateClass(c *game.Class) (err error) {
	return
}

func (csvc *ClassService) GetClass(name string) (class *game.Class, err error) {
	return
}

// GetClasses returns all classes currently in Embercrest
func (csvc *ClassService) GetClasses() (classes []game.Class, err error) {
	return
}
