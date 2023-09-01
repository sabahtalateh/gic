package service

import (
	"fmt"
	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/kind_of_real_project/repo"
)

type Repo interface {
	Select(int) string
}

type Mailing struct {
	userRepo1 Repo
	userRepo2 Repo
}

func (m *Mailing) Send() []string {
	var (
		uu  []string
		out []string
	)

	for i := 0; i < 3; i++ {
		uu = append(uu, m.userRepo1.Select(i))
	}

	for _, u := range uu {
		out = append(out, fmt.Sprintf("sending message to %s", u))
	}

	return out
}

var MailingID = gic.ID("Mailing")

func init() {
	gic.Add[Mailing](
		gic.WithID(MailingID),
		gic.WithInit(func() Mailing {
			return Mailing{
				userRepo1: gic.Get[*repo.UserRepo](gic.WithID(repo.Repo1)),
				userRepo2: gic.Get[*repo.UserRepo](gic.WithID(repo.Repo2)),
			}
		}),
	)
}
