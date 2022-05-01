package znet

import (
	"Hzinx/ziface"
	"fmt"
	"strconv"
)

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	APIs map[uint32]ziface.IRouter
}

// NewMsgHandler 初始化/创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 从Request中找到MsgID
	handler, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("API MsgID =", request.GetMsgID(), "is NOT FOUND! Need register!")
		return
	}
	// 根据MsgID调度对应的router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前Msg绑定的API是否已经存在
	if _, ok := mh.APIs[msgID]; ok {
		panic("repeat API, msgID =" + strconv.Itoa(int(msgID)))
	}
	// 添加Msg与API的绑定关系
	mh.APIs[msgID] = router
	fmt.Println("Add API msgID =", msgID, "succ!")
}
