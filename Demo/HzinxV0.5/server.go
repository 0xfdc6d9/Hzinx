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
	fmt.Println("Call Router Handle")
	fmt.Println("recv from client: msgID =", request.GetMsgID(),
		", data =", string(request.GetMsgData()))

	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	// 创建一个server句柄，使用Hzinx的API
	s := znet.NewServer("[Hzinx V0.5]")

	// 给当前Hzinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	// 启动server
	s.Serve()
}
