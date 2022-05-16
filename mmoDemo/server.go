package main

import (
	"Hzinx/mmoDemo/apis"
	"Hzinx/mmoDemo/core"
	"Hzinx/ziface"
	"Hzinx/znet"
	"fmt"
)

// OnConnectionAdd 当客户端建立连接的时候的hook函数
func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个玩家
	player := core.NewPlayer(conn)

	// 同步当前的PlayerID给客户端，走MsgID:1 消息
	player.SyncPID()

	// 同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()

	// 将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将该连接绑定属性PID
	conn.SetProperty("pID", player.PID)

	// 同步周边玩家上线信息，与现实周边玩家信息
	player.SyncSurrounding()

	fmt.Println("=====> Player pID = ", player.PID, " arrived ====")
}

// OnConnectionLost 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	// 获取当前连接的PID属性
	pID, _ := conn.GetProperty("pID")

	// 根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	// 触发玩家下线业务
	if player != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player ", pID, " left =====")

}

func main() {
	// 创建Hzinx server句柄
	s := znet.NewServer("MMO client Hzinx")
	// 连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	// 启动服务
	s.Serve()
}
