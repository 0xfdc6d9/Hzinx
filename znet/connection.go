package znet

import (
	"Hzinx/utils"
	"Hzinx/ziface"
	"fmt"
	"net"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出的/停止 channel
	ExitChan chan bool

	// 当前连接处理的方法Router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}

	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer c.Stop()
	defer fmt.Println("connID =", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())

	for {
		// 读取客户端的数据到buf中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("c.Conn.Read() occurs an error:", err)
			break
		}

		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			data: buf,
		}
		// 从路由中找到注册绑定的Conn对应的router调用
		c.Router.PreHandle(req)
		c.Router.Handle(req)
		c.Router.PostHandle(req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID =", c.ConnID)

	// 启动从当前连接读数据的业务
	go c.StartReader()
	// TODO 启动从当前连接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()... ConnID =", c.ConnID)

	if c.isClosed { // 如果连接已经关闭
		return
	}
	c.isClosed = true

	c.Conn.Close() // 关闭socket连接
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
