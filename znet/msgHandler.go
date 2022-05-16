package znet

import (
	"Hzinx/utils"
	"Hzinx/ziface"
	"fmt"
	"strconv"
)

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	APIs map[uint32]ziface.IRouter

	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 负责工作worker池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandler 初始化/创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

/*
	采用工作池机制
*/

// StartWorkerPool 启动一个Worker工作池（开启工作池的动作只能发生一次，一个Hzinx框架只能有一个worker工作池）
func (mh *MsgHandler) StartWorkerPool() {
	// 根据WorkerPoolSize分别开启Worker，每个Worker用goroutine来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个Worker被启动
		// 给当前的worker对应的channel开辟空间，第i个worker用第i个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 启动当前的worker，阻塞等待消息从channel中到来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID =", workerID, "is started...")

	// 不断阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息进来，出列的就是一个客户端的request，执行当前request所绑定的业务
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

// SendMsg2TaskQueue 将消息交给TaskQueue，由Worker进行处理
func (mh *MsgHandler) SendMsg2TaskQueue(request ziface.IRequest) {
	// 将消息平均分配给不同的Worker
	// 根据客户端建立的ConnID进行分配（轮询）
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//fmt.Println("Add ConnID =", request.GetConnection().GetConnID(),
	//	"request MsgID =", request.GetMsgID(), "to WorkerID =", workerID)

	// 将消息发送给对应的Worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
