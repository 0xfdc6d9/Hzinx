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

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	// 自定义连接.GetTCPConnection 得到原始的套接字
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..." + "\n"))
	if err != nil {
		fmt.Println("call back before ping error:", err)
		return
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..." + "\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping... error:", err)
		return
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..." + "\n"))
	if err != nil {
		fmt.Println("call back after ping... error:", err)
		return
	}
}

func main() {
	// 创建一个server句柄，使用Hzinx的API
	s := znet.NewServer("[Hzinx V0.3]")

	// 给当前Hzinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 启动server
	s.Serve()
}
