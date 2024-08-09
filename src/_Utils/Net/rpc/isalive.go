package rpc_test

import (
	"net/http"
	"net/rpc"
)

type SyncServer struct {
	From       *Server
	ServerName string
}

func NewSyncServer(from *Server) *SyncServer {
	server := &SyncServer{
		ServerName: "SyncServer",
		From:       from,
	}
	return server

}

func (s *SyncServer) Run(Addr string) {
	rpc.Register(s)
	rpc.RegisterName(s.ServerName, s)
	rpc.HandleHTTP()
	http.ListenAndServe(Addr, nil)
}

type SyncClient struct {
	From       *Server
	ClientName string
	cli        *rpc.Client
	Connected  bool
}

func NewSyncClient(from *Server) *SyncClient {
	cli := &SyncClient{
		ClientName: "SyncClient",
		From:       from,
	}
	return cli
}

// 确认存活方法
type IsAliveArgs struct {
}

type IsAliveReply struct {
	OK bool
}

func (s *SyncServer) IsAlive(args IsAliveArgs, reply *IsAliveReply) error {
	*reply = IsAliveReply{true}
	return nil
}

func (c *SyncClient) SyncIsAlive() bool {
	if !c.Connected {
		err := c.Connect(c.From.SlaveAddr.Address())
		if err != nil {
			c.From.Logger.Printf("SyncIsAlive try to connect error, %v\n", err.Error())
			return false
		}
	}

	rep := IsAliveReply{false}
	err := c.cli.Call("SyncServer.IsAlive", &IsAliveArgs{}, &rep)
	if err != nil {
		c.From.Logger.Printf("RPC cli Call SyncIsAlive error, %v\n", err.Error())
		return false
	}
	return rep.OK
}
