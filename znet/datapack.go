package znet

import (
	"Hzinx/utils"
	"Hzinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

type DataPack struct{}

// NewDataPack 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节）+ ID uint32（4字节）
	return 8
}

// Pack 封包格式：|dataLen|MsgID|data|
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte字节流的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将MsgID写进dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将data数据写进dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// Unpack 将包的Head信息读出来，之后再根据Head信息里的dataLen再进行一次读取
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从[]byte读数据的ioReader
	dataReader := bytes.NewReader(binaryData)

	// 只解压head信息，得到dataLen和MsgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 判断dataLen是否超出MaxPackageSize
	if msg.GetMsgLen() > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data received")
	}

	// 读MsgID
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	return msg, nil
}
