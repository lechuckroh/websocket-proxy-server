package util

import (
	"fmt"
	"rogchap.com/v8go"
)

func AddObjectToGloabl(ctx *v8go.Context, tpl *v8go.ObjectTemplate, name string) error {
	consoleObj, err := tpl.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("failed to create %s instance: %v", name, err)
	}

	global := ctx.Global()
	if err := global.Set(name, consoleObj); err != nil {
		return fmt.Errorf("failed to set %s object to global: %v", name, err)
	}

	return nil
}

type SetObjectFunctions func(fnMap map[string]v8go.FunctionCallback) error

func NewSetObjectFunctions(
	iso *v8go.Isolate,
	objTpl *v8go.ObjectTemplate,
) SetObjectFunctions {
	setObjectFunction := func(name string, fnCallback v8go.FunctionCallback) error {
		errorFnTpl, err := v8go.NewFunctionTemplate(iso, fnCallback)
		if err != nil {
			return fmt.Errorf("failed to create %s FunctionTemplate: %v", name, err)
		}

		if err := objTpl.Set(name, errorFnTpl, v8go.ReadOnly); err != nil {
			return fmt.Errorf("failed to set %s function: %v", name, err)
		}

		return nil
	}

	return func(fnMap map[string]v8go.FunctionCallback) error {
		for name, fn := range fnMap {
			if err := setObjectFunction(name, fn); err != nil {
				return err
			}
		}
		return nil
	}
}
