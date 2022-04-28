package znet

import (
	"Hzinx/ziface"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name      string // 服务器的名称
	IPVersion string // 服务器绑定的IP版本
	IP        string // 服务器监听的IP
	Port      int    // 服务器监听的Port
}

// CallBack2Client 定义当前客户端链接所绑定的handleAPI
// 先写死一个方法，后面再交由用户自定义
func CallBack2Client(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显业务
	fmt.Println("[Conn Handle] CallBack2Client...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("conn.Write() occurs an error:", err)
		return errors.New("CallBack2Client error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	addr, err := net.ResolveTCPAddr("", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.ResolveTCPAddr() occurs an error: ", err)
		return
	}

	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("net.ListenTCP() occurs an error: ", err)
		return
	}
	fmt.Println("start Zinx server  ", s.Name, " success, now listening...")
	var cid uint32

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("listener.AcceptTCP() occurs an error: ", err)
			continue
		}

		// 将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
		dealConn := NewConnection(conn, cid, CallBack2Client)
		cid++

		// 启动当前的链接业务处理
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

/*
初始化Server模块的方法
*/

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}
