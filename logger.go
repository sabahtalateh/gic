package gic

import "go.uber.org/zap"

type withZapSugaredLogger struct {
	s *zap.SugaredLogger
}

func (w withZapSugaredLogger) applyGlobalContainerOption(c *сontainer) {
	c.logger = w.s
}

func WithZapSugaredLogger(s *zap.SugaredLogger) withZapSugaredLogger {
	s = s.With("component", "gic")
	return withZapSugaredLogger{s: s}
}

func (c *сontainer) LogInfof(template string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(template, args...)
	}
}

func (c *сontainer) LogWarnf(template string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Warnf(template, args...)
	}
}
