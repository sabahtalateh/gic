package gic

type id string

type withID struct{ id }

func (w withID) addOption()                   {}
func (w withID) applyGetOption(opts *getOpts) { opts.id = w.id }
func (w withID) applyStageOption(s *stage)    { s.id = w.id }

func WithID(id id) withID {
	return withID{id: id}
}

// ID creates identifier
func ID(value string) id {
	return id(value)
}
