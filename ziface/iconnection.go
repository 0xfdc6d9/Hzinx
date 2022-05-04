package ziface

import "net"

type IConnection interface {
	// Start 启动连接，让当前的连接准备开始工作
	Start()

	// Stop 停止连接 结束当前连接的工作
	Stop()

	// GetTCPConnection 获取当前连接的绑定 socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的 TCP状态 IP:Port
	RemoteAddr() net.Addr

	// SendMsg 提供一个SendMsg方法将我们要发送给客户端的数据，先进行封包，再发送
	SendMsg(msgID uint32, data []byte) error

	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})

	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)

	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}
