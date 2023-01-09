package computer

import (
	"context"
	"fmt"
	"github.com/m4schini/computercraft-go/connection"
	"sync"
	"testing"
)

func NewTestConnection() (out <-chan []byte, in chan<- []byte, conn connection.Connection) {
	outCh := make(chan []byte, 2)
	inCh := make(chan []byte, 2)

	conn = connection.New(inCh, outCh)
	out = outCh
	in = inCh

	return out, in, conn
}

func TestComputer_IsTurtle(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var expected = false
	var actual bool
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		actual, err = c.IsTurtle(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte(fmt.Sprintf("[%v]", expected))
	wg.Wait()

	//assert
	t.Logf("expected: %v", expected)
	t.Logf("  actual: %v", actual)
	if err != nil || actual != expected {
		t.FailNow()
	}
}

func TestComputer_IsPocket(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var expected = true
	var actual bool
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		actual, err = c.IsPocket(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte(fmt.Sprintf("[%v]", expected))
	wg.Wait()

	//assert
	t.Logf("expected: %v", expected)
	t.Logf("  actual: %v", actual)
	if err != nil || actual != expected {
		t.FailNow()
	}
}

func TestComputer_Shutdown(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		err = c.Shutdown(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte("{}")
	wg.Wait()

	//assert
	if err != nil {
		t.FailNow()
	}
}

func TestComputer_Reboot(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		err = c.Reboot(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte("{}")
	wg.Wait()

	//assert
	if err != nil {
		t.FailNow()
	}
}

func TestComputer_Version(t *testing.T) {
	//arrange
	var expected = "CraftOS 1.8"
	var actual string
	var wg sync.WaitGroup
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		actual, err = c.Version(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte(fmt.Sprintf(`["%v"]`, expected))
	wg.Wait()

	//assert
	t.Logf("expected: %v", expected)
	t.Logf("  actual: %v", actual)
	if err != nil || actual != expected {
		t.FailNow()
	}
}

func TestComputer_ComputerId(t *testing.T) {
	//arrange
	var expected = "0"
	var actual string
	var wg sync.WaitGroup
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		actual, err = c.Version(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte(fmt.Sprintf(`["%v"]`, expected))
	wg.Wait()

	//assert
	t.Logf("expected: %v", expected)
	t.Logf("  actual: %v", actual)
	if err != nil || actual != expected {
		t.FailNow()
	}
}

func TestComputer_ComputerLabel(t *testing.T) {
	//arrange
	var expected = "WallE"
	var actual string
	var wg sync.WaitGroup
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		actual, err = c.ComputerLabel(context.TODO())
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte(fmt.Sprintf(`["%v"]`, expected))
	wg.Wait()

	//assert
	t.Logf("expected: %v", expected)
	t.Logf("  actual: %v", actual)
	if err != nil || actual != expected {
		t.FailNow()
	}
}

func TestComputer_SetComputerLabel(t *testing.T) {
	//arrange
	var wg sync.WaitGroup
	var err error
	out, in, conn := NewTestConnection()
	c := NewComputer(conn)

	//act
	go func() {
		wg.Add(1)
		err = c.SetComputerLabel(context.TODO(), "test")
		wg.Done()
	}()

	command := <-out
	t.Logf("command: %v", string(command))
	in <- []byte("{}")
	wg.Wait()

	//assert
	if err != nil {
		t.FailNow()
	}
}
