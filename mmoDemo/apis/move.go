package apis

import (
	"Hzinx/mmoDemo/core"
	"Hzinx/mmoDemo/pb"
	"Hzinx/ziface"
	"Hzinx/znet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传来的proto协议
	msg := &pb.Position{}
	if err := proto.Unmarshal(request.GetMsgData(), msg); err != nil {
		fmt.Println("Move: Position Unmarshal error ", err)
		return
	}
	// 得到当前是哪个玩家在发送位置
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	// 给其他玩家进行当前玩家的位置信息广播（在其他客户端上同步我的移动）
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
