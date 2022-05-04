package znet

import (
	"Hzinx/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	// 管理的连接集合
	connections map[uint32]ziface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}

// Add 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn添加到ConnManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID =", conn.GetConnID(), "added to ConnManager Successfully: conn num =", cm.Len())
}

// Remove 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除连接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID =", conn.GetConnID(), "removed from ConnManager Successfully: conn num =", cm.Len())
}

// Get 根据connID获取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection NOT FOUND")
	}
}

// Len 得到当前连接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清除并终止所有的连接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All connections successfully! conn num = ", cm.Len())
}
