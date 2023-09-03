package internal

import (
	"fmt"
	"github.com/sabahtalateh/gic"
)

type Greeter struct {
	greet string
}

func (g *Greeter) Greet(whom string) string {
	return fmt.Sprintf("%s %s!", g.greet, whom)
}

var RussianGreeter = gic.ID("RussianGreeter")
var ChineseGreeter = gic.ID("ChineseGreeter")

func init() {
	gic.Add[*Greeter](
		gic.WithInit(func() *Greeter { return &Greeter{greet: "Hello"} }),
	)

	gic.Add[*Greeter](
		gic.WithID(RussianGreeter),
		gic.WithInit(func() *Greeter { return &Greeter{greet: "Привет"} }),
	)

	gic.Add[*Greeter](
		gic.WithID(ChineseGreeter),
		gic.WithInit(func() *Greeter { return &Greeter{greet: "你好"} }),
	)
}
