package proxy

import (
	"errors"
	"fmt"
	"github.com/lechuckroh/websocket-proxy-server/polyfill/util"
	"rogchap.com/v8go"
)

func InjectTo(ctx *v8go.Context) (Proxy, error) {
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	iso, err := ctx.Isolate()
	if err != nil {
		return nil, fmt.Errorf("failed to get v8go.Isolate: %v", err)
	}

	proxy := NewProxy()

	proxyTpl, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy ObjectTemplate: %v", err)
	}

	setObjectFunctions := util.NewSetObjectFunctions(iso, proxyTpl)

	if err := setObjectFunctions(map[string]v8go.FunctionCallback{
		"onInit":                       proxy.GetOnInitFunctionCallback(),
		"onDestroy":                    proxy.GetOnDestroyFunctionCallback(),
		"addReceivedMessageMiddleware": proxy.GetAddReceivedMessageMiddlewareFunctionCallback(),
		"addSentMessageMiddleware":     proxy.GetAddSentMessageMiddlewareFunctionCallback(),
		"addReceiveMessageMiddleware":  proxy.GetAddReceiveMessageMiddlewareFunctionCallback(),
		"addSendMessageMiddleware":     proxy.GetAddSendMessageMiddlewareFunctionCallback(),
	}); err != nil {
		return nil, err
	}

	return proxy, util.AddObjectToGloabl(ctx, proxyTpl, "proxy")
}
