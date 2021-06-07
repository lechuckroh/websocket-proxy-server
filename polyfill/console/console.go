package console

import (
	"fmt"
	"io"
	"log"
	"os"
	"rogchap.com/v8go"
)

type Console interface {
	GetErrorFunctionCallback() v8go.FunctionCallback
	GetLogFunctionCallback() v8go.FunctionCallback
}

type consoleImpl struct {
	ErrorWriter io.Writer
	LogWriter   io.Writer
}

func NewConsole(opts ...Option) Console {
	console := &consoleImpl{
		ErrorWriter: os.Stderr,
		LogWriter:   os.Stdout,
	}

	for _, opt := range opts {
		opt.apply(console)
	}

	return console
}

func (c *consoleImpl) GetErrorFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		arguments := make([]interface{}, len(info.Args()))
		for i, arg := range info.Args() {
			arguments[i] = arg
		}

		_, err := fmt.Fprintln(c.ErrorWriter, arguments...)
		if err != nil {
			log.Printf("[ERR] console.error: %v", err)
		}

		return nil
	}
}

func (c *consoleImpl) GetLogFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		arguments := make([]interface{}, len(info.Args()))
		for i, arg := range info.Args() {
			arguments[i] = arg
		}
		_, err := fmt.Fprintln(c.LogWriter, arguments...)
		if err != nil {
			log.Printf("[ERR] console.log: %v", err)
		}

		return nil
	}
}
