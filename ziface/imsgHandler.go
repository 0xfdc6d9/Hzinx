package ziface

type IMsgHandler interface {
	// DoMsgHandler 调度/执行对应的 router 消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
