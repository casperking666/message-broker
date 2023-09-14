package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/distributed3/pairbroker/stubs"
)

func getOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()
	return localAddr
}

func multiplyTwoNumbers(x int, y int) int {
	return x * y
}

type Factory struct{}

//TODO: Define a Multiply function to be accessed via RPC.
func (f *Factory) Multiply(req stubs.Pair, res *stubs.JobReport) (err error) {
	value := multiplyTwoNumbers(req.X, req.Y)
	res.Result = value
	res.Message = "job run successfully"
	return
}

func main() {
	pAddr := flag.String("ip", "8050", "IP and port to listen on")
	brokerAddr := flag.String("broker", "127.0.0.1:8030", "Address of broker instance")
	flag.Parse()
	//TODO: You'll need to set up the RPC server, and subscribe to the running broker instance.
	server, _ := rpc.Dial("tcp", *brokerAddr)
	callbackFuncStr := "Factory.Multiply"
	ip := getOutboundIP()
	fmt.Println(ip)
	rpc.Register(&Factory{})
	listener, err := net.Listen("tcp", ":"+"8050") // this order doesn't matter fuck, fucking retarded gpt
	defer listener.Close()

	subsription := stubs.Subscription{Topic: "multiply", FactoryAddress: ip + ":" + *pAddr, Callback: callbackFuncStr}
	status := new(stubs.StatusReport)
	err = server.Call(stubs.Subscribe, subsription, status)

	if err != nil {
		fmt.Println("sth wrong for some reasons")
		fmt.Println(err)
	}

	if err != nil {
		fmt.Println("server reporst")
		fmt.Println(err)
	}
	rpc.Accept(listener)
}
