package main

import "Hzinx/znet"

func main() {
	s := znet.NewServer("[Hzinx V0.1]")
	s.Serve()
}
