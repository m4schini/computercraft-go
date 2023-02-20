package connection

import (
	"context"
	"fmt"
	"github.com/m4schini/logger"
	"sync"
	"testing"
)

func TestConn_Execute(t *testing.T) {
	//arrange
	var expected = `test`
	var actual []interface{}
	var wg sync.WaitGroup
	var err error
	in := make(chan []byte)
	out := make(chan []byte)
	conn := New(in, out, WithLog(logger.Named("test")))

	//act
	wg.Add(1)
	go func() {
		actual, err = conn.Execute(context.TODO(), "test")
		wg.Done()
	}()

	outgoing := <-out
	t.Logf("<- outgoing: %v", string(outgoing))
	incoming := fmt.Sprintf(`["%v"]`, expected)
	t.Logf("-> incoming: %v", incoming)
	in <- []byte(incoming)
	wg.Wait()

	//assert
	e := make([]interface{}, 1)
	e[0] = string("test")
	t.Logf("expected: %v", e)
	t.Logf("  actual: %v", actual)
	if err != nil || actual[0] != e[0] {
		t.FailNow()
	}
}

func TestConn_Execute_TimedOut(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var err error
	in := make(chan []byte)
	out := make(chan []byte)
	conn := New(in, out)
	ctx, cancel := context.WithCancel(context.TODO())

	//act
	wg.Add(1)
	go func() {
		_, err = conn.Execute(ctx, "test")
		wg.Done()
	}()

	outgoing := <-out
	t.Logf("<- outgoing: %v", string(outgoing))
	cancel()
	wg.Wait()

	//assert
	t.Logf("expected: %v", context.Canceled)
	t.Logf("  actual: %v", err)
	if err == nil || err != context.Canceled {
		t.FailNow()
	}
}

func TestConn_Execute_closedIn(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var err error
	in := make(chan []byte)
	out := make(chan []byte)
	conn := New(in, out)

	//act
	close(in)

	wg.Add(1)
	go func() {
		_, err = conn.Execute(context.TODO(), "test")
		wg.Done()
	}()

	outgoing := <-out
	t.Logf("<- outgoing: %v", string(outgoing))
	wg.Wait()

	//assert
	t.Logf("Error: %v", err)
	if err == nil {
		t.FailNow()
	}
}

func TestConn_Execute_closedOut(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var err error
	in := make(chan []byte)
	out := make(chan []byte)
	conn := New(in, out)

	//act
	close(out)

	wg.Add(1)
	go func() {
		_, err = conn.Execute(context.TODO(), "test")
		wg.Done()
	}()
	wg.Wait()

	//assert
	t.Logf("Error: %v", err)
	if err == nil {
		t.FailNow()
	}
}

func TestSendOnClosedChannel(t *testing.T) {
	var err any
	func() {
		ch := make(chan string, 1)
		close(ch)
		defer func() {
			x := recover()
			err = x
		}()
		ch <- "1"
	}()

	t.Logf("Err: %v (%T)", err, err)
}
