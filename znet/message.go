package znet

type Message struct {
	ID      uint32
	Data    []byte
	DataLen uint32
}

// NewMsgPackage 创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(u uint32) {
	m.ID = u
}

func (m *Message) SetData(bytes []byte) {
	m.Data = bytes
}

func (m *Message) SetDataLen(u uint32) {
	m.DataLen = u
}
