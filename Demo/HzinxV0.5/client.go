package main

import (
	"Hzinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client start...")
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("net.Dial() occurs an error: ", err)
		return
	}

	for {
		// 发送封包的Message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("HzinxV0.5 client test message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 服务器给客户端回复一个Message
		// 先读取流中的headData，得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read headData error:", err)
			break
		}
		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 再根据dataLen读取data
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("---> Recv Server Msg : ID =", msg.GetMsgId(), ", len =", msg.GetMsgLen(), ", data =", string(msg.GetData()))
		}

		// CPU阻塞
		time.Sleep(1 * time.Second)
	}
}
