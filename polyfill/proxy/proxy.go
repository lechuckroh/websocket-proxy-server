package proxy

import (
	"log"
	"rogchap.com/v8go"
)

type ExecuteMiddlewaresFn func(*v8go.Value, ...*v8go.Value) (*v8go.Value, error)

type Proxy interface {
	GetOnInitFunctionCallback() v8go.FunctionCallback
	GetOnDestroyFunctionCallback() v8go.FunctionCallback

	GetAddReceivedMessageMiddlewareFunctionCallback() v8go.FunctionCallback
	GetAddSentMessageMiddlewareFunctionCallback() v8go.FunctionCallback

	GetAddResponseToBackendMessageMiddlewareFunctionCallback() v8go.FunctionCallback
	GetAddResponseToClientMessageMiddlewareFunctionCallback() v8go.FunctionCallback

	ExecuteReceivedMessageMiddlewares(*v8go.Value, ...*v8go.Value) (*v8go.Value, error)
	ExecuteSentMessageMiddlewares(*v8go.Value, ...*v8go.Value) (*v8go.Value, error)
	ExecuteResponseToBackendMessageMiddlewares(*v8go.Value, ...*v8go.Value) (*v8go.Value, error)
	ExecuteResponseToClientMessageMiddlewares(*v8go.Value, ...*v8go.Value) (*v8go.Value, error)

	ExecuteOnInit() error
	ExecuteOnDestroy() error
}

type proxyImpl struct {
	initFn    *v8go.Function
	destroyFn *v8go.Function

	receivedMessageMiddlewares []*v8go.Function
	sentMessageMiddlewares     []*v8go.Function

	responseToBackendMessageMiddlewares []*v8go.Function
	responseToClientMessageMiddlewares  []*v8go.Function
}

func NewProxy() Proxy {
	proxy := proxyImpl{}

	return &proxy
}

func (p *proxyImpl) GetOnInitFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) != 1 {
			log.Printf("[ERR] onInit() requires 1 argument")
			return nil
		}

		onInitFn, err := args[0].AsFunction()
		if err != nil {
			log.Printf("[ERR] onInit() argument is not function: %v", err)
			return nil
		}

		p.initFn = onInitFn
		return nil
	}
}

func (p *proxyImpl) GetOnDestroyFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) != 1 {
			log.Printf("[ERR] onDestroy() requires 1 argument")
			return nil
		}

		onDestroyFn, err := args[0].AsFunction()
		if err != nil {
			log.Printf("[ERR] onDestroy() argument is not function: %v", err)
			return nil
		}

		p.destroyFn = onDestroyFn
		return nil
	}
}

func (p *proxyImpl) GetAddReceivedMessageMiddlewareFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		for i, arg := range args {
			middleware, err := arg.AsFunction()
			if err != nil {
				log.Printf("[ERR] addReceivedMessageMiddleware() args[%d] is not function: %v", i, err)
				continue
			}
			p.receivedMessageMiddlewares = append(p.receivedMessageMiddlewares, middleware)
		}

		return nil
	}
}

func (p *proxyImpl) GetAddSentMessageMiddlewareFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		for i, arg := range args {
			middleware, err := arg.AsFunction()
			if err != nil {
				log.Printf("[ERR] addSentMessageMiddleware() args[%d] is not function: %v", i, err)
				continue
			}
			p.sentMessageMiddlewares = append(p.sentMessageMiddlewares, middleware)
		}

		return nil
	}
}

func (p *proxyImpl) GetAddResponseToBackendMessageMiddlewareFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		for i, arg := range args {
			middleware, err := arg.AsFunction()
			if err != nil {
				log.Printf("[ERR] addResponseToBackendMessageMiddleware() args[%d] is not function: %v", i, err)
				continue
			}
			p.responseToBackendMessageMiddlewares = append(p.responseToBackendMessageMiddlewares, middleware)
		}

		return nil
	}
}

func (p *proxyImpl) GetAddResponseToClientMessageMiddlewareFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		for i, arg := range args {
			middleware, err := arg.AsFunction()
			if err != nil {
				log.Printf("[ERR] addResponseToClientMessageMiddleware() args[%d] is not function: %v", i, err)
				continue
			}
			p.responseToClientMessageMiddlewares = append(p.responseToClientMessageMiddlewares, middleware)
		}

		return nil
	}
}

func (p *proxyImpl) ExecuteReceivedMessageMiddlewares(message *v8go.Value, rest ...*v8go.Value) (*v8go.Value, error) {
	return newExecuteMiddlewaresFunc(p.receivedMessageMiddlewares)(message, rest...)
}

func (p *proxyImpl) ExecuteSentMessageMiddlewares(message *v8go.Value, rest ...*v8go.Value) (*v8go.Value, error) {
	return newExecuteMiddlewaresFunc(p.sentMessageMiddlewares)(message, rest...)
}

func (p *proxyImpl) ExecuteResponseToBackendMessageMiddlewares(message *v8go.Value, rest ...*v8go.Value) (*v8go.Value, error) {
	return newExecuteMiddlewaresFunc(p.responseToBackendMessageMiddlewares)(message, rest...)
}

func (p *proxyImpl) ExecuteResponseToClientMessageMiddlewares(message *v8go.Value, rest ...*v8go.Value) (*v8go.Value, error) {
	return newExecuteMiddlewaresFunc(p.responseToClientMessageMiddlewares)(message, rest...)
}

func newExecuteMiddlewaresFunc(middlewares []*v8go.Function) ExecuteMiddlewaresFn {
	return func(message *v8go.Value, rest ...*v8go.Value) (*v8go.Value, error) {
		prevMessage := message
		for _, middleware := range middlewares {
			args := []v8go.Valuer{prevMessage}
			for _, arg := range rest {
				args = append(args, arg)
			}

			result, err := middleware.Call(args...)
			if err != nil {
				return nil, err
			}
			if result == nil || result.IsNullOrUndefined() {
				return nil, nil
			}
			prevMessage = result
		}
		return prevMessage, nil
	}
}

func (p *proxyImpl) ExecuteOnInit() error {
	if p.initFn != nil {
		_, err := p.initFn.Call()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *proxyImpl) ExecuteOnDestroy() error {
	if p.destroyFn != nil {
		_, err := p.destroyFn.Call()
		if err != nil {
			return err
		}
	}
	return nil
}
