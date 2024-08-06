package rpccommon

type Args struct {
	Arg1 int
	Arg2 int
}

type Reply struct {
	Sum int
}

type RPCServer struct {
}

func (s *RPCServer) Add(args Args, reply *int) error {
	*reply = args.Arg1 + args.Arg2
	return nil
}
