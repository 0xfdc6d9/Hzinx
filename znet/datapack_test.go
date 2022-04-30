package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// dataPack拆包封包的单元测试

func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	// 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	// 创建一个go负责从客户端处理业务
	go func() {
		// 从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
				continue
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// ------------> 拆包的过程 <---------------
				// 定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					// 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					// io.ReadFull 从conn读数据直至注满headData
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("io.ReadFull() occurs an error:", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg中有数据，需要进行第二次读取
						// 第二次从conn读，根据head中的dataLen，再读取data内容

						// 想把data字段添加到msgHead之后，需要进行类型断言（将接口转为具体的类型），将msgHead转换成Message结构（向下转换）
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据dataLen再次从io流中读取数据
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Println("==> Recv Msg: ID =", msg.ID, ", len =", msg.DataLen, ", data =", string(msg.Data))
					}

				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial, err:", err)
		return
	}

	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		ID:      1,
		Data:    []byte{'H', 'z', 'i', 'n', 'x'},
		DataLen: 5,
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	// 封装地热个msg2包
	msg2 := &Message{
		ID:      2,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
		DataLen: 7,
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	_, err = conn.Write(sendData1)
	if err != nil {
		return
	}

	// 客户端阻塞，等待服务端返回
	select {}
}
