package gic

// Init runs all the added components init function (see: Add).
func Init() error {
	globC.mu.Lock()
	defer globC.mu.Unlock()

	if globC.initialized {
		return ErrInitialized
	}

	for i, fn := range globC.initFns {
		if _, ok := globC.initsDone[i]; ok {
			continue
		}
		if err := fn(globC); err != nil {
			return err
		}
		globC.initsDone[i] = struct{}{}
	}

	if globC.dump != nil {
		writeDump(globC)
	}

	globC.initialized = true

	return nil
}
