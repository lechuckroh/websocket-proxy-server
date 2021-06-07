package timer

import (
	"bytes"
	"fmt"
	"github.com/lechuckroh/websocket-proxy-server/polyfill/console"
	"rogchap.com/v8go"
	"strings"
	"testing"
	"time"
)

func TestInjectTo(t *testing.T) {
	iso, _ := v8go.NewIsolate()

	global, _ := v8go.NewObjectTemplate(iso)
	if err := InjectTo(iso, global); err != nil {
		t.Error(err)
	}

	// inject console object for testing
	ctx, _ := v8go.NewContext(iso, global)
	logBuf := new(bytes.Buffer)
	if err := console.InjectTo(ctx, console.WithLog(logBuf)); err != nil {
		t.Error(err)
	}

	// run setInterval()
	startScript := `
		setInterval(function(arg1, arg2) {
			console.log(arg1 + " " + arg2 + " " + Date.now());
		}, 400, "1", "2");
	`
	var intervalID int32
	if result, err := ctx.RunScript(startScript, ""); err != nil {
		t.Error(err)
	} else {
		intervalID = result.Int32()
	}

	time.Sleep(time.Second)

	// run clearInterval()
	clearScript := fmt.Sprintf("clearInterval(%d);", intervalID)
	if _, err := ctx.RunScript(clearScript, ""); err != nil {
		t.Error(err)
	}

	// assert log messages
	logMessages := strings.Split(logBuf.String(), "\n")
	executionCount := 0
	for _, logMessage := range logMessages {
		if logMessage != "" {
			executionCount++

			if !strings.HasPrefix(logMessage, "1 2 ") {
				t.Errorf("restArgs are not passed. message: %s", logMessage)
			}
		}
	}

	if executionCount != 2 {
		t.Errorf("execution count mismatch. 2 != %d", executionCount)
	}
}
