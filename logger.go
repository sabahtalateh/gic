package gic

import "go.uber.org/zap"

type withZapSugaredLogger struct {
	s *zap.SugaredLogger
}

func (w withZapSugaredLogger) applyGlobalContainerOption(c *container) {
	c.logger = w.s
}

func WithZapSugaredLogger(s *zap.SugaredLogger) withZapSugaredLogger {
	s = s.With("component", "gic")
	return withZapSugaredLogger{s: s}
}

func (c *container) LogInfof(template string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(template, args...)
	}
}

func (c *container) LogWarnf(template string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Warnf(template, args...)
	}
}
