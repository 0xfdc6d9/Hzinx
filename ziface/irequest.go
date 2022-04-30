package ziface

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetMsgData GetData 得到请求的消息数据
	GetMsgData() []byte
	GetMsgID() uint32
}
