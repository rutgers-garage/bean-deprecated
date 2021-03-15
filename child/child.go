package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

func main() {

	listenParent()
	for {

	}
}

func isOnline() bool {
	return true
}

func execWebService(endpoint string) {
	cmd := exec.Command("dir")

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(stdout))
}

func execShellService() {
	// TODO
}

// listen for parent connection
func listenParent() {
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Parent")

	fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

	go handleConnection(c)
}
func handleConnection(conn net.Conn) {
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			fmt.Println("Client left.")
			conn.Close()
			return
		}

		log.Println("Client message:", string(buffer[:len(buffer)-1]))

		conn.Write(buffer)
		mes := string(buffer)

		// parse the data
		req := strings.Split(mes, " ")

		if req[0] == "web" {
			execWebService(req[1])
		}

	}
}
