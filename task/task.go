package task

import (
	"context"
	"sync"
	"time"
)

type State uint8

const (
	Ready State = iota
	Processing
	Completed
	Failed
)

type Step func() error

type task struct {
	state State
	steps []Step
	error error

	mu sync.Mutex

	_subscribersMu sync.Mutex
	_subscribers   map[chan struct{}]struct{}
}

func (t *task) Start() {
	go func() {
		if t.state != Ready {
			return
		}
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.steps == nil {
			t.steps = make([]Step, 0)
		}

		var err error
		t.state = Processing
		for _, step := range t.steps {
			err = step()
			if err != nil {
				break
			}
		}
		if err != nil {
			t.error = err
			t.state = Failed
		} else {
			t.state = Completed
		}
		t.notify()
	}()
}

func (t *task) State() State {
	return t.state
}

func (t *task) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (t *task) Done() <-chan struct{} {
	t._subscribersMu.Lock()
	if t._subscribers == nil {
		t._subscribers = make(map[chan struct{}]struct{})
	}
	t._subscribersMu.Unlock()

	ch := make(chan struct{}, 2)
	t._subscribers[ch] = struct{}{}

	if t.state == Completed || t.state == Failed {
		ch <- struct{}{}
	}
	return ch
}

func (t *task) notify() {
	if t._subscribers == nil {
		return
	}

	for s := range t._subscribers {
		select {
		case s <- struct{}{}:
			break
		default:
			break
		}
	}
}

func (t *task) Err() error {
	return t.error
}

func (t *task) Value(key any) any {
	return nil
}

type Task interface {
	Start()
	State() State

	context.Context
}
