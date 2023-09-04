package gic

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type getOpts struct{ id id }
type getOpt interface{ applyGetOption(*getOpts) }

// Get return component from container
// errors: ErrNotFound
func Get[T any](opts ...getOpt) (T, error) {
	oo := getOpts{}
	for _, opt := range opts {
		opt.applyGetOption(&oo)
	}

	return get[T](globC, oo)
}

// MustGet return component from container. panics on Get error
// errors: ErrNotFound
func MustGet[T any](opts ...getOpt) T {
	t, err := Get[T](opts...)
	check(err)
	return t
}

// errors: ErrNotFound
func get[T any](c *container, opts getOpts) (t T, err error) {
	lookFor := reflect.TypeOf(&t).Elem()
	if err = checkGetType(lookFor); err != nil {
		return t, err
	}

	// T found
	if comps, ok := c.components[lookFor]; ok {
		comp, err := compByID(comps, lookFor, opts.id)
		if err != nil {
			return t, err
		}
		return comp.c.(T), nil
	}

	// Try to find type assignable to T
	for typ, comps := range c.components {
		if !typ.AssignableTo(lookFor) {
			continue
		}

		comp, err := compByID(comps, typ, opts.id)
		if err != nil {
			return t, err
		}

		return reflect.ValueOf(comp.c).Convert(lookFor).Interface().(T), nil
	}

	// T not found. Try hint
	if err = hint(c.components, lookFor, makeCaller()); err != nil {
		return t, errors.Join(ErrNotFound, err)
	}

	// T Not found
	return t, errors.Join(ErrNotFound, fmt.Errorf("%s[id=%s] not found %s", lookFor, opts.id, makeCaller()))
}

func compByID(comps map[string]*component, t reflect.Type, id id) (*component, error) {
	comp, ok := comps[id.v]
	if !ok {
		var (
			foundMsg string
			found    = compsForErr(comps, t)
		)
		if len(found) > 0 {
			foundMsg = "\nexisting components:\n" + strings.Join(found, "\n")
		}
		return nil, errors.Join(ErrNotFound, fmt.Errorf("%s[id=%s] not found\n%s%s", t, id, makeCaller(), foundMsg))
	}
	return comp, nil
}

func compsForErr(comps map[string]*component, t reflect.Type) []string {
	// sort keys by caller
	keys := maps.Keys(comps)
	sort.Slice(keys, func(i, j int) bool { return comps[keys[i]].caller.String() < comps[keys[j]].caller.String() })

	found := make([]string, len(keys))
	for i, k := range keys {
		v := comps[k]
		found[i] = fmt.Sprintf("%s[id=%s] at %s", t, k, v.caller)
	}

	return found
}

func checkGetType(t reflect.Type) error {
	if t.Kind() == reflect.Interface {
		return errors.Join(
			ErrInterface,
			fmt.Errorf("attempting to get %s interface\n%s", t, makeCaller()),
		)
	}

	return nil
}
