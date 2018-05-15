package genserver

// Server rocks
type Server struct {
	Name     string
	Callback GenServer
	state    State

	reqCh   chan GenReq
	replyCh chan Reply
}

// GenReq combines call, cast and shutdown requests
type GenReq struct {
	Msg
	t int // call, cast or shutdown
}

// Msg of request
type Msg interface{}

// Reply from server
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
	s.reqCh = make(chan GenReq)
	s.replyCh = make(chan Reply)
	s.state = s.Callback.Init()
	go s.serve()
}

func (s *Server) serve() {
	for {
		req := <-s.reqCh
		switch req.t {
		case call:
			reply := s.Callback.HandleCall(req, &(s.state))
			s.replyCh <- reply
		case cast:
			s.Callback.HandleCast(req, &(s.state))
		case shutdown:
			close(s.reqCh)
			close(s.replyCh)
			s.Callback.Terminate(s.state)
			return
		}
	}
}

// Call req
func (s *Server) Call(msg Msg) (reply Reply) {
	s.reqCh <- GenReq{Msg: msg, t: call}
	reply = <-s.replyCh
	return
}

// Cast req
func (s *Server) Cast(msg Msg) {
	s.reqCh <- GenReq{Msg: msg, t: cast}
}

// Stop req
func (s *Server) Stop() {
	s.reqCh <- GenReq{t: shutdown}
}
