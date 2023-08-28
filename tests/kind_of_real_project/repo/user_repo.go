package repo

import (
	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/kind_of_real_project/system"
)

type UserRepo struct {
	db *system.DB
}

func (u *UserRepo) Select(id int) string {
	switch id {
	case 0:
		return "Ivan"
	case 1:
		return "Petr"
	default:
		return "Anonymous"
	}
}

var Repo1, Repo2 = gic.ID("UserRepo1"), gic.ID("UserRepo2")

// var UserRepo1ID = ic.id()
// var UserRepo2ID = ic.id()

func init() {
	gic.Add[*UserRepo](
		gic.WithID(Repo1),
		gic.WithInit(func() *UserRepo { return &UserRepo{db: gic.Get[*system.DB]()} }),
	)

	gic.Add[*UserRepo](
		gic.WithID(Repo2),
		gic.WithInit(func() *UserRepo { return &UserRepo{db: gic.Get[*system.DB]()} }),
	)
}
