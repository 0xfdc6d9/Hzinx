package znet

import (
	"Hzinx/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name      string // 服务器的名称
	IPVersion string // 服务器绑定的IP版本
	IP        string // 服务器监听的IP
	Port      int    // 服务器监听的Port
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept() occurs an error: ", err)
			continue
		}

		go func() {
			for {

				buf := make([]byte, 1024)
				cnt, err := conn.Read(buf)
				if err != nil {
					fmt.Println("conn.Read() occurs an error: ", err)
					continue
				}
				if cnt == 0 {
					continue
				}

				if _, err := conn.Write(buf); err != nil {
					fmt.Println("conn.Write() occurs an error: ", err)
					continue
				}
			}
		}()
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
