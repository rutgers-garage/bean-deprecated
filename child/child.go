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

	listenParent("192.168.25.183", "8080")
}

func isOnline() bool {
	return true
}

func execWebService(endpoint string) {
	cmd := exec.Command("chromium-browser", endpoint)

	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(stdout))
}

// TODO: Make it so it kills specific web apps instead of all of them?
func killWebService() {
	cmd := exec.Command("pkill", "-o", "chromium")
	cmd.Start()

}

func execShellService() {
	// TODO
}

// listen for parent connection
func listenParent(ip string, port string) {
	fmt.Println("Listening on: " + ip + ":" + port)
	l, err := net.Listen("tcp", ip+":"+port)
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

	handleConnection(c)
}
func handleConnection(conn net.Conn) {
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			fmt.Println("Client left.")
			conn.Close()
			return
		}

		mes := string(buffer[:len(buffer)-1])

		log.Println("Client message:", mes)

		conn.Write(buffer)

		// parse the data
		req := strings.Split(mes, " ")

		if req[0] == "web" {
			execWebService(req[1])
		} else if req[0] == "terminate" {
			killWebService()
		}

	}
}
