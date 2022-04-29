package main

import "Hzinx/znet"

func main() {
	s := znet.NewServer("[Hzinx V0.2]")
	s.Serve()
}
