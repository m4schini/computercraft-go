package task

import (
	"testing"
	"time"
)

func TestTask_Done(t *testing.T) {
	tsk := new(task)

	tsk.steps = []Step{
		func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		},
	}

	start := time.Now()
	t.Log("ready", time.Since(start))
	tsk.Start()
	t.Log("processing...", time.Since(start))
	<-tsk.Done()
	t.Log("done", time.Since(start))
	<-tsk.Done()
	t.Log("done", time.Since(start))

}
