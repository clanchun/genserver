package genserver

// Server rocks
type Server struct {
	name     string
	state    State
	reqCh    chan GenReq
	callback GenServer
}

// GenReq combines call, cast and shutdown requests
type GenReq struct {
	msg
	ch chan Reply
	t  int // call, cast or shutdown
}

// request msg
type msg interface{}

// Reply of server
type Reply interface{}

// State of server
type State interface{}

const (
	call = iota
	cast
	shutdown
)

// Start server
func (s *Server) Start() {
	reqCh := make(chan GenReq)
	s.reqCh = reqCh
	s.state = s.callback.Init()
	go s.serve(reqCh)
	return
}

func (s *Server) serve(reqCh chan GenReq) {
	for {
		req := <-reqCh
		switch req.t {
		case call:
			reply := s.callback.HandleCall(req, &(s.state))
			req.ch <- reply
			close(req.ch)
		case cast:
			s.callback.HandleCast(req, &(s.state))
		case shutdown:
			s.callback.Terminate(s.state)
			return
		}
	}
}

// Call req
func (s *Server) Call(msg msg) (reply Reply) {
	replyCh := make(chan Reply)
	s.reqCh <- GenReq{msg: msg, ch: replyCh, t: call}
	reply = <-replyCh
	return
}

// Cast req
func (s *Server) Cast(msg msg) {
	s.reqCh <- GenReq{msg: msg, t: cast}
}

// Stop req
func (s *Server) Stop() {
	s.reqCh <- GenReq{t: shutdown}
}
