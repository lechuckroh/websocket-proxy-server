package console

import (
	"errors"
	"fmt"
	"github.com/lechuckroh/websocket-proxy-server/polyfill/util"
	"rogchap.com/v8go"
)

// InjectTo injects console object to v8go.Context.
func InjectTo(ctx *v8go.Context, opts ...Option) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return fmt.Errorf("failed to get v8go.Isolate: %v", err)
	}

	console := NewConsole(opts...)

	consoleTpl, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return fmt.Errorf("failed to create console ObjectTemplate: %v", err)
	}

	setObjectFunctions := util.NewSetObjectFunctions(iso, consoleTpl)

	if err := setObjectFunctions(map[string]v8go.FunctionCallback{
		"error": console.GetErrorFunctionCallback(),
		"log":   console.GetLogFunctionCallback(),
	}); err != nil {
		return err
	}

	return util.AddObjectToGloabl(ctx, consoleTpl, "console")
}
