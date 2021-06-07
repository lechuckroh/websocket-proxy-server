package timer

import (
	"errors"
	"log"
	"rogchap.com/v8go"
)

type Timers interface {
	GetSetIntervalFunctionCallback() v8go.FunctionCallback
	GetClearIntervalFunctionCallback() v8go.FunctionCallback

	GetSetTimeoutFunctionCallback() v8go.FunctionCallback
	GetClearTimeoutFunctionCallback() v8go.FunctionCallback
}

type timersImpl struct {
	lastTimerID  int32
	timerTaskMap map[int32]*timerTask
}

func NewTimers() Timers {
	return &timersImpl{
		lastTimerID:  0,
		timerTaskMap: make(map[int32]*timerTask),
	}
}

func (t *timersImpl) nextTimerID() int32 {
	t.lastTimerID++
	return t.lastTimerID
}

func (t *timersImpl) GetSetIntervalFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		args := info.Args()

		timerID, err := t.startTimerTask(args, true)
		if err != nil {
			log.Printf("[ERR] failed to start setInterval function: %v", err)
			return int32Value(ctx, 0)
		}

		return int32Value(ctx, timerID)
	}
}

func (t *timersImpl) GetClearIntervalFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) == 1 && args[0].IsInt32() {
			intervalID := args[0].Int32()
			t.clearTimerTask(intervalID, true)
		}
		return nil
	}
}

func (t *timersImpl) GetSetTimeoutFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		ctx := info.Context()
		args := info.Args()

		timerID, err := t.startTimerTask(args, false)
		if err != nil {
			log.Printf("[ERR] failed to start setTimeout function: %v", err)
			return int32Value(ctx, 0)
		}

		return int32Value(ctx, timerID)
	}
}

func (t *timersImpl) GetClearTimeoutFunctionCallback() v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) == 1 && args[0].IsInt32() {
			timeoutID := args[0].Int32()
			t.clearTimerTask(timeoutID, false)
		}
		return nil
	}
}

func (t *timersImpl) startTimerTask(args []*v8go.Value, interval bool) (int32, error) {
	argsCount := len(args)

	if argsCount == 0 {
		return 0, errors.New("at least 1 argument required")
	}

	// 1st argument: function
	fn, err := args[0].AsFunction()
	if err != nil {
		return 0, err
	}

	// 2nd argument: delay
	delay := int32(0)
	if argsCount >= 2 {
		if args[1].IsInt32() {
			delay = args[1].Int32()
		}
	}

	// rest arguments
	restArgs := make([]v8go.Valuer, 0)
	if argsCount > 2 {
		for _, arg := range args[2:] {
			restArgs = append(restArgs, arg)
		}
	}

	// create timer task
	timerID := t.nextTimerID()
	task := &timerTask{
		timerID:   timerID,
		delay:     delay,
		interval:  interval,
		arguments: restArgs,

		timerCallback: func(args ...v8go.Valuer) {
			_, err := fn.Call(args...)
			if err != nil {
				log.Printf("[ERR] failed to call timer function: %v", err)
			}
		},
		clearCallback: func(timerID int32) {
			delete(t.timerTaskMap, timerID)
		},
	}

	t.timerTaskMap[timerID] = task

	// start timer task
	task.Start()

	return timerID, nil
}

func (t *timersImpl) clearTimerTask(timerID int32, interval bool) {
	if timerID <= 0 {
		return
	}
	if task, ok := t.timerTaskMap[timerID]; ok && task.interval == interval {
		task.Clear()
	}
}

func int32Value(ctx *v8go.Context, i32 int32) *v8go.Value {
	iso, _ := ctx.Isolate()
	value, _ := v8go.NewValue(iso, i32)
	return value
}
