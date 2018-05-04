package genserver

// GenServer any specific server needs to implement this interface
type GenServer interface {
	Init() State
	HandleCall(GenReq, *State) Reply
	HandleCast(GenReq, *State)
	Terminate(State)
}
