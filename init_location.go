package gic

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"reflect"
	"sort"
	"strings"
)

type initLocation struct {
	fn     int // init function index
	caller caller
}

func (c *container) setInitLocation(t reflect.Type, id string, initIndex int, caller caller) {
	if _, ok := c.initLocs[t]; !ok {
		c.initLocs[t] = map[string]initLocation{}
	}
	c.initLocs[t][id] = initLocation{fn: initIndex, caller: caller}
}

func (c *container) initLocation(t reflect.Type, id string) (initLocation, error) {
	ids, ok := c.initLocs[t]
	if !ok {
		// T not found. Try hint
		return initLocation{}, errors.Join(ErrNotFound, hintType(c.initLocs, t, makeCaller()))
	}

	loc, ok := ids[id]
	if !ok {
		var (
			foundMsg string
			found    = initLocsForErr(ids, t)
		)
		if len(found) > 0 {
			foundMsg = "\nexisting components:\n" + strings.Join(found, "\n")
		}

		return initLocation{}, errors.Join(
			ErrNotFound, fmt.Errorf("%s%s not found\n%s%s", t, strID(id), makeCaller(), foundMsg),
		)
	}

	return loc, nil
}

func initLocsForErr(locs map[string]initLocation, t reflect.Type) []string {
	keys := maps.Keys(locs)
	sort.Slice(keys, func(i, j int) bool { return locs[keys[i]].caller.String() < locs[keys[j]].caller.String() })

	found := make([]string, len(keys))
	for i, k := range keys {
		v := locs[k]
		found[i] = fmt.Sprintf("%s%s at %s", t, strID(k), v.caller)
	}

	return nil
}
