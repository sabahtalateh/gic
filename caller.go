package gic

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type caller struct {
	found    bool
	funcName string
	file     string
	line     int
}

func (c caller) String() string {
	if !c.found {
		return "[caller not found]"
	}
	return fmt.Sprintf("%s\n\t%s:%d", c.funcName, c.file, c.line)
}

var skipCallers = []string{
	"github.com/sabahtalateh/gic.",
	"runtime",
}

func checkCallFromInit(c caller) error {
	skip := 1

	for {
		pc, _, _, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		fName := runtime.FuncForPC(pc).Name()

		parts := strings.Split(fName, ".")
		if len(parts) < 2 {
			skip++
			continue
		}
		last := parts[len(parts)-1]
		preLast := parts[len(parts)-2]
		if _, err := strconv.ParseInt(last, 10, 64); err == nil && preLast == "init" {
			return nil
		}

		skip++
	}

	return errors.Join(ErrNotFromInit, fmt.Errorf("%s\n", c))
}

// returns first call made outside github.com/sabahtalateh/gic
func makeCaller() caller {
	// skip self
	skip := 1

FOR:
	for {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}

		fName := runtime.FuncForPC(pc).Name()
		for _, excl := range skipCallers {
			if strings.HasPrefix(fName, excl) {
				skip++
				continue FOR
			}
		}

		return caller{found: true, funcName: fName, file: file, line: line}
	}

	return caller{found: false}
}
