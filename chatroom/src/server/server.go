package serve

import (
	IMsystem "IM_system/src/proto"
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type RPCServer struct {
}

func (s *RPCServer) SendTo(ctx context.Context, req *IMsystem.SendToReq) (*IMsystem.SendToResp, error) {
	return nil, nil
}
func (s *RPCServer) Broadcast(ctx context.Context, req *IMsystem.BroadcastReq) (*IMsystem.BroadcastResp, error) {
	return nil, nil
}
func (s *RPCServer) ChangeName(ctx context.Context, req *IMsystem.ChangeNameReq) (*IMsystem.ChangeNameResp, error) {
	return nil, nil
}

type Server struct {
	IP               string
	Port             int
	RPCPort          int
	OnlineUsers      map[string]*User
	Lock             sync.Mutex
	BroadcastChannel chan string
}

func (s *Server) SendTo(ctx context.Context, req *IMsystem.SendToReq) (*IMsystem.SendToResp, error) {
	return nil, nil
}
func (s *Server) Broadcast(ctx context.Context, req *IMsystem.BroadcastReq) (*IMsystem.BroadcastResp, error) {
	return nil, nil
}
func (s *Server) ChangeName(ctx context.Context, req *IMsystem.ChangeNameReq) (*IMsystem.ChangeNameResp, error) {
	newName := req.NewName

	_, hasBennUsed := s.OnlineUsers[newName]
	u := s.OnlineUsers[req.Name]
	if hasBennUsed {
		u.SendMessage("this name has been used...\n")
		return &IMsystem.ChangeNameResp{
			Ok:         false,
			ErrMessage: "name has been used",
		}, fmt.Errorf("name has been used")
	}
	u.server.Lock.Lock()

	delete(u.server.OnlineUsers, req.Name)
	u.server.OnlineUsers[newName] = u

	u.server.Lock.Unlock()
	u.Name = newName
	u.SendMessage("your name update success : " + u.Name + "\n")
	return &IMsystem.ChangeNameResp{
		Ok:         true,
		ErrMessage: "",
	}, nil
}

// NewServer 创建一个Server
func NewServer(ip string, port, rpcport int) *Server {
	server := &Server{
		IP:               ip,
		Port:             port,
		RPCPort:          rpcport,
		OnlineUsers:      make(map[string]*User),
		BroadcastChannel: make(chan string),
	}

	return server
}

func (thisServer *Server) Broadcast_(user *User, msg string) {
	broadMsg := "[" + user.IP + "]" + user.Name + ":" + msg + "\n"
	thisServer.BroadcastChannel <- broadMsg
}

func (thisServer *Server) HandleBroadcast() {
	for {
		msg := <-thisServer.BroadcastChannel
		thisServer.Lock.Lock()

		for _, user := range thisServer.OnlineUsers {
			// user.MessageChannel <- msg
			user.SendMessage(msg)
		}

		thisServer.Lock.Unlock()
	}
}

// HandleUser 处理一个user的业务
func (thisServer *Server) HandleUser(conn net.Conn) {
	user := NewUser(conn, thisServer)

	user.Online()
	defer user.Offline()

	isAlive := make(chan bool)
	buffer := make([]byte, 4096)
	for {
		// 用户保活检测
		select {
		case <-isAlive:
			// break
			// 不做任何事，但是刷新time，表明用户连接活动
		case <-time.After(3600 * time.Second):
			user.SendMessage("you are over_time and kicked")
			// close(user.MessageChannel)
			err := conn.Close()
			if err != nil {
				fmt.Println("Server.HandleUser conn.Close error: ", err)
			}
			// break
			return
		}

		// 用户业务处理
		sz, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Server.HandleUser conn.Read error: ", err)
			return
		} else if sz == 0 {
			user.Offline()
			return
		} else {
			msg := string(buffer[:sz-1])
			user.HandleMessage(msg)
			isAlive <- true
		}
	}
}

func (thisServer *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", thisServer.IP, thisServer.Port))
	if err != nil {
		fmt.Println("Server.Start net.Listen error: ", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Server.Start, net.Close error : ", err)
			return
		}
	}(listener)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", thisServer.IP, thisServer.RPCPort))
	if err != nil {
		fmt.Println("network error", err)
	}
	defer ln.Close()

	//创建grpc服务
	srv := grpc.NewServer()
	//注册服务
	// IMsystem.RegisterIMServiceServer(srv, &thisServer{})
	go srv.Serve(ln)
	// if err != nil {
	// 	fmt.Println("Serve error", err)
	// }
	defer srv.Stop()
	fmt.Println("go srv.Serve(ln)")

	go thisServer.HandleBroadcast()
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Server.Start, listener.Accept error: ", err)
			continue
		}
		fmt.Printf("a new connection from %v\n", conn.RemoteAddr().String())

		// do handler
		go thisServer.HandleUser(conn)
	}
}

func MainServer() {
	server := NewServer("127.0.0.1", 8888, 8889)
	server.Start()
}
