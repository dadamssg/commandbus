package commandbus

import (
	"bytes"
	"testing"
)

type TestCommand struct{}

type TestCommand2 struct{}

func TestNew(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Log("New command bus not created!")
		t.Fail()
	}
}

func TestCanGetRegisteredHandlerFunc(t *testing.T) {
	bus := New()
	test := 0
	handler1 := func(cmd interface{}) { test = 1 }
	handler2 := func(cmd interface{}) { test = 2 }

	bus.RegisterHandler(&TestCommand{}, handler1)
	bus.RegisterHandler(&TestCommand2{}, handler2)

	bus.GetHandler(&TestCommand2{})(nil)

	if test != 2 {
		t.Log("Wrong handler called!")
		t.Fail()
	}

	bus.GetHandler(&TestCommand{})(nil)

	if test != 1 {
		t.Log("Wrong handler called!")
		t.Fail()
	}
}

func TestCanUseHandler(t *testing.T) {
	var buffer bytes.Buffer
	bus := New()
	command := &TestCommand{}
	bus.RegisterHandler(command, func(cmd interface{}) {
		buffer.WriteString("executed")
	})

	bus.Handle(command)

	if buffer.String() != "executed" {
		t.Log("Command was not executed!")
		t.Fail()
	}
}

func TestCanUseMiddleware(t *testing.T) {
	var buffer bytes.Buffer
	bus := New()
	command := &TestCommand{}
	bus.RegisterHandler(command, func(cmd interface{}) {
		buffer.WriteString("executed")
	})

	bus.AddMiddleware(0, func(cmd interface{}, next HandlerFunc) {
		buffer.WriteString("0")
		next(cmd)
		buffer.WriteString("0")
	})

	bus.Handle(command)

	if buffer.String() != "0executed0" {
		t.Log("Command was not executed!")
		t.Fail()
	}
}

func TestCanUsePrioritizedMiddleware(t *testing.T) {
	var buffer bytes.Buffer
	bus := New()
	command := &TestCommand{}
	bus.RegisterHandler(command, func(cmd interface{}) {
		buffer.WriteString("a")
	})

	bus.AddMiddleware(0, func(cmd interface{}, next HandlerFunc) {
		buffer.WriteString("0")
		next(cmd)
		buffer.WriteString("0")
	})

	bus.AddMiddleware(1, func(cmd interface{}, next HandlerFunc) {
		buffer.WriteString("1")
		next(cmd)
		buffer.WriteString("1")
	})

	bus.Handle(command)

	if buffer.String() != "10a01" {
		t.Log("Execution occurred out of order!")
		t.Fail()
	}
}
