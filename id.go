package gic

import "fmt"

type id string

type withID struct {
	id id
	c  caller
}

func (w withID) addOption()                   {}
func (w withID) addOptionCallInfo() string    { return fmt.Sprintf("gic.WithID\n%s", w.c) }
func (w withID) applyGetOption(opts *getOpts) { opts.id = w.id }
func (w withID) applyStageOption(s *stage)    { s.id = w.id }

func WithID(id id) withID {
	return withID{id: id, c: makeCaller()}
}

// ID creates identifier
func ID(value string) id {
	return id(value)
}
