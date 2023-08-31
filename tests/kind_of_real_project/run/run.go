package run

import (
	"context"
	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/kind_of_real_project/service"
	"go.uber.org/zap"
)

const dump = "/Users/kravtsov777/Code/go/src/github.com/sabahtalateh/gic/tests/kind_of_real_project/dump"

func Run() []string {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	_ = gic.ConfigureGlobalContainer(
		gic.WithZapSugaredLogger(sugar),
		gic.WithDump(
			gic.WithDumpDir(dump),
			// gic.WithOverrideRoot(
			// 	"/Users/kravtsov777/Code/go/src/github.com/sabahtalateh/gic/tests/kind_of_real_project",
			// 	"/Users/kravtsov777/Code/go/src/github.com/sabahtalateh/gic/tests/kind_of_real_project2",
			// ),
		),
	)

	if err := gic.Init(); err != nil {
		panic(err)
	}

	if err := gic.Start(context.Background()); err != nil {
		panic(err)
	}

	m, err := gic.GetE[*service.Mailing](gic.WithID(service.MailingID))
	if err != nil {
		panic(err)
	}

	res := m.Send()

	if err := gic.Stop(context.Background()); err != nil {
		panic(err)
	}

	return res
}
