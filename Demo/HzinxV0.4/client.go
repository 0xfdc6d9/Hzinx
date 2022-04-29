package main

import (
	"Hzinx/utils"
	"fmt"
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
		_, err := conn.Write([]byte("Hello Hzinx V0.2..."))
		if err != nil {
			fmt.Println("conn.Write() occurs an error: ", err)
			continue
		}

		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read() occurs an error: ", err)
			continue
		}
		if cnt == 0 {
			continue
		}
		fmt.Printf("Server calls back: %s, cnt = %d\n", buf, cnt)

		time.Sleep(1 * time.Second)
	}
}
