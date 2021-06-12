package proxy

import (
	"bytes"
	"github.com/lechuckroh/websocket-proxy-server/polyfill/console"
	"rogchap.com/v8go"
	"strings"
	"testing"
)

func TestInjectTo(t *testing.T) {
	iso, _ := v8go.NewIsolate()
	global, _ := v8go.NewObjectTemplate(iso)
	ctx, _ := v8go.NewContext(iso, global)

	// inject proxy
	proxy, err := InjectTo(ctx)
	if err != nil {
		t.Error(err)
	}

	// inject console object for testing
	logBuf := new(bytes.Buffer)
	if err := console.InjectTo(ctx, console.WithLog(logBuf)); err != nil {
		t.Error(err)
	}

	// run script
	script := `
		proxy.onInit(function() { console.log('onInit'); });
		proxy.onDestroy(function() { console.log('onDestroy'); });
		proxy.addReceivedMessageMiddleware(
			function(msg) { 
				console.log('receivedMessage1 ' + msg);
				return msg + '1';
			},
			function(msg) { 
				console.log('receivedMessage2 ' + msg);
				return msg + '2';
			},
		);
		proxy.addSentMessageMiddleware(
			function(msg) { 
				console.log('sentMessage1 ' + msg);
				return msg + '1';
			},
			function(msg) { 
				console.log('sentMessage2 ' + msg);
				return msg + '2';
			},
		);
		proxy.addResponseToBackendMessageMiddleware(
			function(msg) { 
				console.log('responseToBackend1 ' + msg);
				return msg + '1';
			},
			function(msg) { 
				console.log('responseToBackend2 ' + msg);
				return msg + '2';
			},
		);
		proxy.addResponseToClientMessageMiddleware(
			function(msg) { 
				console.log('responseToClient1 ' + msg);
				return msg + '1';
			},
			function(msg) { 
				console.log('responseToClient2 ' + msg);
				return msg + '2';
			},
		);
	`
	if _, err := ctx.RunScript(script, ""); err != nil {
		t.Error(err)
	}

	// call proxy functions
	msg, _ := v8go.NewValue(iso, "foo")

	_ = proxy.ExecuteOnInit()
	_, _ = proxy.ExecuteResponseToBackendMessageMiddlewares(msg)
	_, _ = proxy.ExecuteResponseToClientMessageMiddlewares(msg)
	_, _ = proxy.ExecuteReceivedMessageMiddlewares(msg)
	_, _ = proxy.ExecuteSentMessageMiddlewares(msg)
	_ = proxy.ExecuteOnDestroy()

	// assert log messages
	expectedLogMessage := strings.Join([]string{
		"onInit",
		"responseToBackend1 foo",
		"responseToBackend2 foo1",
		"responseToClient1 foo",
		"responseToClient2 foo1",
		"receivedMessage1 foo",
		"receivedMessage2 foo1",
		"sentMessage1 foo",
		"sentMessage2 foo1",
		"onDestroy",
		"",
	}, "\n")
	logMessage := logBuf.String()

	if logMessage != expectedLogMessage {
		t.Errorf("actual logMessages: %s", logMessage)
	}
}
