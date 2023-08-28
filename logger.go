package gic

import "go.uber.org/zap"

type withZapSugaredLogger struct {
	s *zap.SugaredLogger
}

func (w withZapSugaredLogger) applyGlobalContainerOption(c *Container) {
	c.logger = w.s
}

func WithZapSugaredLogger(s *zap.SugaredLogger) withZapSugaredLogger {
	s = s.With("component", "gic")
	return withZapSugaredLogger{s: s}
}

func (c *Container) LogInfof(template string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(template, args...)
	}
}

func (c *Container) LogWarnf(template string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Warnf(template, args...)
	}
}
