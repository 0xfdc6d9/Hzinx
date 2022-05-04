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

// DoConnBegin 创建连接之后执行的Hook函数
func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("=====> DoConnBegin is called...")
	if err := conn.SendMsg(202, []byte("DoConn BEGIN")); err != nil {
		fmt.Println(err)
	}

	// 给当前的连接设置一些属性
	fmt.Println("Set conn Name, Age...")
	conn.SetProperty("Name", "Nepenthe8")
	conn.SetProperty("Age", 22)
}

// DoConnEnd 连接销毁前执行的Hook函数
func DoConnEnd(conn ziface.IConnection) {
	fmt.Println("=====> DoConnEnd is called...")
	fmt.Println("connID =", conn.GetConnID(), "is lost...")

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name:", name)
	}
	if age, err := conn.GetProperty("Age"); err == nil {
		fmt.Println("Age", age)
	}
}

func main() {
	// 创建一个server句柄，使用Hzinx的API
	s := znet.NewServer("[Hzinx V0.9]")

	// 注册连接的Hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnEnd)

	// 给当前Hzinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloHzinxRouter{})

	// 启动server
	s.Serve()
}
