package znet

import (
	"Hzinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前的连接状态
	isClosed bool

	// 告知当前连接已经退出的/停止 channel（由Reader告知Writer退出）
	ExitChan chan bool

	// 消息管理的MsgID 和对应的处理业务的API
	MsgHandler ziface.IMsgHandler

	// 无缓冲管道，用于读、写Goroutine之间通信
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}

	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer c.Stop()
	defer fmt.Println("[Reader is exit], connID =", c.ConnID, "remote addr is", c.RemoteAddr().String())

	for {
		// 读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("c.Conn.Read() occurs an error:", err)
		//	break
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的msgHead
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}

		// 拆包，得到msgID和dataLen放在一个msg对象中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}

		// 根据dataLen，再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中找到注册绑定的Conn对应的router调用
		// 根据绑定好的MsgID找到对应的API执行
		go c.MsgHandler.DoMsgHandler(req)
	}
}

// SendMsg 提供一个SendMsg方法将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	// 将data进行封包 DataLen|ID|Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id =", msgID)
		return errors.New("pack error msg")
	}

	// 将数据交给Writer发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

// StartWriter 给客户端发送消息
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	// 循环等待channel消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也需要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID =", c.ConnID)

	// 启动从当前连接读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()... ConnID =", c.ConnID)

	if c.isClosed { // 如果连接已经关闭
		return
	}
	c.isClosed = true

	if err := c.Conn.Close(); err != nil {
		return
	} // 关闭socket连接

	// 告知Writer关闭
	c.ExitChan <- true

	close(c.msgChan)
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
