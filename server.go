package hup

import (
	"net"
)

type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//开启业务服务方法
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(route string, msgHandler HandlerFunc)
}

type Server struct {
	//tcp4 or other
	IPVersion string
	//服务绑定的地址
	Addr string

	msgHandler map[string]HandlerFunc
}

func NewServer(addr string, IPVersion string) *Server {
	return &Server{
		IPVersion:  IPVersion,
		Addr:       addr,
		msgHandler: make(map[string]HandlerFunc),
	}
}

func (s *Server) AddRoute(route string, msgHandlerFunc HandlerFunc) {
	s.msgHandler[route] = msgHandlerFunc
}

func (s *Server) Start() {
	InfoF("Listener at %s is starting\n", s.Addr)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, s.Addr)
		if err != nil {
			Info("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			Info("listen", s.IPVersion, "err", err)
			return
		}
		Info("Server listening")

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				Error("accept err:", err)
				return
			}
			Info("get conn remote addr:", conn.RemoteAddr().String())
			dealConn := NewConnection(conn, s.msgHandler)

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}

	}()
}

func (s *Server) Stop() {
	//panic("implement me")
	s.Start()

	select {}
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

type HandlerFunc func(*Connection, *Message)
