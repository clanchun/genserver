package genserver

import (
	"fmt"
	"testing"
	"time"
)

var callback GenServer

type counterserver string

func (c counterserver) Init() State {
	return 0
}

func (c counterserver) HandleCast(req GenReq, state *State) {
	switch req.msg {
	case "incr":
		*state = (*state).(int) + 1
	}
}

func (c counterserver) HandleCall(req GenReq, state *State) Reply {
	var n int
	switch req.msg {
	case "get":
		n = (*state).(int)
	}
	return n
}

func (c counterserver) Terminate(state State) {
	fmt.Println("counter server terminated.")
	return
}

func TestMain(t *testing.T) {
	var callback GenServer
	var c counterserver = "counter"
	callback = c
	var s = &Server{name: "counter", callback: callback}

	s.Start()
	fmt.Println("counter server started.")

	n := s.Call("get")
	fmt.Println("counter is: ", s.state)
	if n != 0 {
		t.Error("init test failed")
	}

	s.Cast("incr")

	n = s.Call("get")
	fmt.Println("counter is: ", s.state)
	if n != 1 {
		t.Error("call test failed")
	}

	s.Stop()

	time.Sleep(time.Second)
}
