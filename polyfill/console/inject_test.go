package console

import (
	"bytes"
	"rogchap.com/v8go"
	"testing"
)

func TestInjectTo(t *testing.T) {
	iso, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(iso)

	errorBuf := new(bytes.Buffer)
	logBuf := new(bytes.Buffer)
	if err := InjectTo(ctx, WithError(errorBuf), WithLog(logBuf)); err != nil {
		t.Error(err)
	}

	script := `
		console.error('foo');
		console.log('bar');
	`
	if _, err := ctx.RunScript(script, ""); err != nil {
		t.Error(err)
	}

	actualError := errorBuf.String()
	actualLog := logBuf.String()
	expectedError := "foo\n"
	expectedLog := "bar\n"

	if actualError != expectedError {
		t.Errorf("error output: '%s' != '%s'", actualError, expectedError)
	}
	if actualLog != expectedLog {
		t.Errorf("log output: '%s' != '%s'", actualLog, expectedLog)
	}
}
