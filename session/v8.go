package session

import (
	"github.com/lechuckroh/websocket-proxy-server/polyfill/console"
	"github.com/lechuckroh/websocket-proxy-server/polyfill/timer"
	"rogchap.com/v8go"
)

// initV8 initialize v8go
func initV8() (*v8go.Context, error) {
	if iso, err := v8go.NewIsolate(); err != nil {
		return nil, err
	} else if global, err := v8go.NewObjectTemplate(iso); err != nil {
		return nil, err
	} else if ctx, err := v8go.NewContext(iso, global); err != nil {
		return nil, err
	} else {
		if err := timer.InjectTo(iso, global); err != nil {
			return nil, err
		} else if err := console.InjectTo(ctx); err != nil {
			return nil, err
		} else {
			return ctx, nil
		}
	}
}
