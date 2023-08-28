package gic

// Init runs all the added (see: Add) components init function
func Init() error {
	globC.mu.Lock()
	defer globC.mu.Unlock()

	if globC.initialized {
		return ErrInitialized
	}

	for _, fn := range globC.initFns {
		if err := fn(globC); err != nil {
			return err
		}
	}

	if globC.dump != nil {
		writeDump(globC)
	}

	globC.initialized = true

	return nil
}
