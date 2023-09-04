package run

import (
	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/kind_of_real_project/service"
)

type A struct {
}

func init() {
	gic.Add[[]A](
		gic.WithInit(func() []A {
			return []A{{}, {}}
		}),
	)
}

func Run() []string {
	if err := gic.Init(); err != nil {
		panic(err)
	}

	m, err := gic.Get[service.Mailing](gic.WithID(service.MailingID))
	if err != nil {
		panic(err)
	}

	res := m.Send()

	return res
}
