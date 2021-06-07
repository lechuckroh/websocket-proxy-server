package timer

import (
	"rogchap.com/v8go"
	"time"
)

type timerTask struct {
	timerID   int32
	delay     int32
	interval  bool
	cleared   bool
	finished  bool
	arguments []v8go.Valuer

	timerCallback func(args ...v8go.Valuer)
	clearCallback func(int32)
}

func (t *timerTask) Start() {
	go func() {
		// clear on timer finished
		defer t.Clear()

		ticker := time.NewTicker(time.Duration(t.delay) * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			if t.cleared || t.finished {
				break
			}

			if t.timerCallback != nil {
				t.timerCallback(t.arguments...)
			}

			if !t.interval {
				t.finished = true
				break
			}
		}
	}()
}

func (t *timerTask) Clear() {
	t.cleared = true

	if t.clearCallback != nil {
		t.clearCallback(t.timerID)
	}
}
