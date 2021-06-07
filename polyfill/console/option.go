package console

import "io"

type Option interface {
	apply(c *consoleImpl)
}

type optionFunc func(c *consoleImpl)

func (f optionFunc) apply(c *consoleImpl) {
	f(c)
}

func WithLog(output io.Writer) Option {
	return optionFunc(func(c *consoleImpl) {
		c.LogWriter = output
	})
}

func WithError(output io.Writer) Option {
	return optionFunc(func(c *consoleImpl) {
		c.ErrorWriter = output
	})
}
