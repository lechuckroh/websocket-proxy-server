package timer

import (
	"github.com/lechuckroh/websocket-proxy-server/polyfill/util"
	"rogchap.com/v8go"
)

func InjectTo(iso *v8go.Isolate, global *v8go.ObjectTemplate) error {
	timer := NewTimers()

	objFnsAdder := util.NewSetObjectFunctions(iso, global)

	return objFnsAdder(map[string]v8go.FunctionCallback{
		"clearInterval": timer.GetClearIntervalFunctionCallback(),
		"setInterval":   timer.GetSetIntervalFunctionCallback(),
		"clearTimeout":  timer.GetClearTimeoutFunctionCallback(),
		"setTimeout":    timer.GetSetTimeoutFunctionCallback(),
	})
}
