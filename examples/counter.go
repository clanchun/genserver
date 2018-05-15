package main

import (
	"fmt"
	"time"

	"github.com/clanchun/genserver"
)

type counter string

const (
	get = iota
	incr
	decr
)

func (c counter) Init() genserver.State {
	return 0
}

func (c counter) HandleCast(req genserver.GenReq, state *genserver.State) {
	switch req.Msg {
	case incr:
		*state = (*state).(int) + 1
	case decr:
		*state = (*state).(int) - 1
	}
}

func (c counter) HandleCall(req genserver.GenReq, state *genserver.State) genserver.Reply {
	var n int
	switch req.Msg {
	case get:
		n = (*state).(int)
	}
	return n
}

func (c counter) Terminate(state genserver.State) {
	fmt.Println("counter server terminated.")
	return
}

func main() {
	var c counter = "counter"
	var s = &genserver.Server{Name: "counter", Callback: c}

	s.Start()
	fmt.Println("counter server started.")

	n := s.Call(get)
	fmt.Println("counter is:", n)

	s.Cast(incr)

	n = s.Call(get)
	fmt.Println("counter is:", n)

	s.Cast(decr)

	n = s.Call(get)
	fmt.Println("counter is:", n)

	s.Stop()

	time.Sleep(time.Second)
}
