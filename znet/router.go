package znet

import "Hzinx/ziface"

// BaseRouter 实现IRouter接口
// 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写就好了（不需要三个方法都重写）
type BaseRouter struct{}

// 之所以BaseRouter的方法都为空
// 是因为有的Router不希望有PreHandle，PostHandle这两个业务
// 所以Router全部继承BaseRouter的好处就是，可以不实现PreHandle，PostHandle

func (BaseRouter) PreHandle(request ziface.IRequest) {}

func (BaseRouter) Handle(request ziface.IRequest) {}

func (BaseRouter) PostHandle(request ziface.IRequest) {}
