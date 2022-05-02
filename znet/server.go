package znet

import (
	"Hzinx/utils"
	"Hzinx/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name       string // 服务器的名称
	IPVersion  string // 服务器绑定的IP版本
	IP         string // 服务器监听的IP
	Port       int    // 服务器监听的Port
	MsgHandler ziface.IMsgHandler
}

func (s *Server) Start() {
	fmt.Printf("[Hzinx] Server name: %s, listener at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	fmt.Printf("[Hzinx] Version %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	// 开启消息队列及Worker工作池
	s.MsgHandler.StartWorkerPool()

	addr, err := net.ResolveTCPAddr("", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.ResolveTCPAddr() occurs an error:", err)
		return
	}

	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("net.ListenTCP() occurs an error:", err)
		return
	}
	fmt.Println("start Hzinx server", s.Name, "success, now listening...")
	var cid uint32

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("listener.AcceptTCP() occurs an error:", err)
			continue
		}

		// 将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
		dealConn := NewConnection(conn, cid, s.MsgHandler)
		cid++

		// 启动当前的连接业务处理
		go dealConn.Start()
	}
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	go s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!")
}

/*
初始化Server模块的方法
*/

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandler(),
	}

	return s
}
