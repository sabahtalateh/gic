package system

import (
	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/kind_of_real_project/config"
)

type DB struct {
	dsn string
}

func (d *DB) Query(q string) string {
	if q == "select user where id = 1" {
		return "Ivan"
	}

	if q == "select user where id = 2" {
		return "Petr"
	}

	if q == "select user where id = 3" {
		return "Vasisualiy"
	}

	return "anonymous"
}

func init() {
	gic.Add[*DB](
		gic.WithInit(func() *DB {
			return &DB{dsn: gic.Get[*config.Config]().DB.DSN}
		}),
	)
}
