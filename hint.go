package gic

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// https://github.com/golang/go/blob/release-branch.go1.21/src/reflect/type.go#L160
func hasElem(t reflect.Type) bool {
	return slices.Contains([]reflect.Kind{
		reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice,
	}, t.Kind())
}

func hint(comps map[reflect.Type]map[id]*component, lookFor reflect.Type, c *caller) error {
	for typ := range comps {
		if hasElem(typ) {
			if typ.Elem().AssignableTo(lookFor) {
				return fmt.Errorf("type %s not found but %s found\n%s", lookFor, typ, c)
			}
		}

		if hasElem(lookFor) {
			if lookFor.Elem().AssignableTo(typ) {
				return fmt.Errorf("type %s not found but %s found\n%s", lookFor, typ, c)
			}
		}
	}
	return nil
}
