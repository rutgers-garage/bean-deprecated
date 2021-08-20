package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os/exec"
)

type Child struct{}

func (c *Child) ExecWebService(endpoint *string, childRetString *string) error {
	cmd := exec.Command("explorer", *endpoint)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

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
		log.Fatal(fmt.Printf("Unable to listen on given port: 8080", err))
	}

	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}

// func main() {

// 	listenParent("", "8080")
// }


// TODO: Make it so it kills specific web apps instead of all of them?
func killWebService() {
	cmd := exec.Command("pkill", "-o", "chromium")
	cmd.Run()
}

func execShellService() {
	// TODO
}

// listen for parent connection
// func listenParent(ip , port ) {
// 	fmt.Println("Listening on: " + ip + ":" + port)
// 	l, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer l.Close()
// 	for {
// 		c, err := l.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Println("Connected to Parent")

// 		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

// 		handleConnection(c)
// 	}
// }

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	fmt.Println("Waiting for parent message")
	buffer, err := reader.ReadBytes('\n')
	fmt.Println("After buffer read")
	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	mes := (buffer[:len(buffer)-1])

	log.Println("Client message:", mes)

	conn.Write(buffer)

	// parse the data
	// req := strings.Split(mes, " ")

	// if req[0] == "web" {
	// 	go execWebService(req[1])
	// } else if req[0] == "terminate" {
	// 	fmt.Println("Terminating chromium")
	// 	killWebService()
	// } else if req[0] == "poll" {
	// 	conn.Close()
	// 	return
	// }
}
