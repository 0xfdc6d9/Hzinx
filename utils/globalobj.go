package utils

import (
	"Hzinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
存储一切有关Hzinx框架的全局参数，供其他模块使用
一些参数是可以通过Hzinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TCPServer ziface.IServer // 当前Hzinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TCPPort   int            // 当前服务器主机监听的Port
	Name      string         // 当前服务器的名称
	/*
		Hzinx
	*/
	Version        string // 当前Hzinx的版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 // 当前Hzinx框架数据包的最大值
}

// GlobalObject 定义一个全局的对外GlobalOb对象
var GlobalObject *GlobalObj

// Reload 从Hzinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/Hzinx.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}

// 当import utils包之后，先调utils包中的init方法。重复导入只会执行一次

// 初始化GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		TCPServer:      nil,
		Host:           "0.0.0.0",
		TCPPort:        8999,
		Name:           "HzinxServerAPP",
		Version:        "V0.5",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// 尝试从conf/Hzinx.json 加载用户自定义的参数
	GlobalObject.Reload()
}
