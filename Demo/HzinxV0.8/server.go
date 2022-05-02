package main

import (
	"Hzinx/ziface"
	"Hzinx/znet"
	"fmt"
)

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter // 在一个结构体中包含一个实现了接口的结构体
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	fmt.Println("recv from client: msgID =", request.GetMsgID(),
		", data =", string(request.GetMsgData()))

	if err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
		return
	}
}

// HelloHzinxRouter Hello Hzinx 自定义路由
type HelloHzinxRouter struct {
	znet.BaseRouter
}

func (h *HelloHzinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloHzinxRouter Handle...")
	fmt.Println("recv from client: msgID =", request.GetMsgID(),
		", data =", string(request.GetMsgData()))

	if err := request.GetConnection().SendMsg(201, []byte("Hello! Welcome to Hzinx!")); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	// 创建一个server句柄，使用Hzinx的API
	s := znet.NewServer("[Hzinx V0.5]")

	// 给当前Hzinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloHzinxRouter{})

	// 启动server
	s.Serve()
}
