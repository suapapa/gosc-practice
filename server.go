package main

import (
	"./gosc"
	"fmt"
	"net"
	"os"
)

func messageHander(m *osc.Message) {
	fmt.Println(m)
}

func runOSCServer() {
	srv := &osc.UDPServer{
		Address:    ":8000",
		Dispatcher: osc.NewDumpDispatcher(os.Stdout),
	}

	fmt.Println(net.ResolveUDPAddr("udp", srv.Address))
	srv.ListenAndServe()
}
