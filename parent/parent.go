package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type PollStatus struct {
	MachineName string
	Status      bool
}

type Service struct {
	Title       string
	Endpoint    string
	ServiceType string
}

type Machine struct {
	Ip        string
	Supported []Service
}

var machines map[string]Machine

func main() {
	// Find all the machines from json
	parseJSON()

	// start
	connectChild("192.168.25.183", "8080")

}

func connectChild(ip string, port string) {
	fmt.Println("Trying to connect to child. " + ip + ":" + port)
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection established")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Text to send: ")

		input, _ := reader.ReadString('\n')

		fmt.Fprint(conn, input)

		message, _ := bufio.NewReader(conn).ReadString('\n')

		log.Print("Server relay:", message)
	}

}

// check which machines are online
// TODO: set the status correctly
func pollMachines() []PollStatus {
	statuses := make([]PollStatus, 0, len(machines))
	for k := range machines {
		statuses = append(statuses, PollStatus{MachineName: k, Status: true})
	}
	return statuses
}

// executes this command on specified machine
func executeCommand(machineName string, cmd string) bool {

	statuses := pollMachines()
	statusFlag := false

	// check if machine is available
	for _, machine := range statuses {
		if machine.MachineName == machineName {
			statusFlag = true
		}
	}

	if !statusFlag {
		fmt.Println("Machine is not available")
		return false
	}

	// check if machine has support for this command
	if !checkCommand(machineName, cmd) {
		fmt.Println(fmt.Sprintf("Machine does not support %s", cmd))
		return false
	}

	return true
}

// check if command can work on specfied machine
func checkCommand(machineName string, cmd string) bool {

	for _, supported := range machines[machineName].Supported {
		if supported.Title == cmd {
			return true
		}
	}
	return false
}

func parseJSON() {
	jsonFile, err := os.Open("services.json")
	if err != nil {
		fmt.Println(err)
	}

	byteVal, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteVal, &machines)
}
