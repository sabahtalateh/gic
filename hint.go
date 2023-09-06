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

func hintType(comps map[reflect.Type]map[string]initLocation, t reflect.Type, c caller) error {
	for typ := range comps {
		if hasElem(typ) {
			if typ.Elem().AssignableTo(t) {
				return fmt.Errorf("type %s not found but %s found\n%s", t, typ, c)
			}
		}

		if hasElem(t) {
			if t.Elem().AssignableTo(typ) {
				return fmt.Errorf("type %s not found but %s found\n%s", t, typ, c)
			}
		}
	}
	return nil
}
