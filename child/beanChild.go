package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os/exec"
)

type Child struct{}

func (c *Child) ExecWebService(endpoint *string, childRetString *string) error {
	cmd := exec.Command("explorer", *endpoint)

	stdout, err := cmd.Output()
	if err != nil {
		log.Println(*endpoint, err)
	}

	log.Println(*endpoint, stdout)

	*childRetString = "200"

	return nil
	// fmt.Println((stdout))
}

func main() {
	c := &Child{}
	rpc.Register(c)

	fmt.Println("CHILD UP ON 8080")
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(fmt.Println("Unable to listen on given port: 8080", err))
	}

	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}