package run

import (
	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/kind_of_real_project/service"
)

func Run() []string {
	if err := gic.Init(); err != nil {
		panic(err)
	}

	m, err := gic.GetE[service.Mailing](gic.WithID(service.MailingID))
	if err != nil {
		panic(err)
	}

	res := m.Send()

	return res
}
