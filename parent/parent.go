package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

	// set up http server
	http.HandleFunc("/bean", httpStuff)
	http.ListenAndServe(":8081", nil)
}

func httpStuff(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/bean" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	fmt.Println("Connection established")
	switch req.Method {
	case "GET":
		fmt.Println("GET request")
		http.ServeFile(w, req, "test.html")
	case "POST":
		fmt.Println("POST request")
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		command := req.FormValue("command")
		machine := req.FormValue("machine")
		fmt.Fprintf(w, "Machine = %s\n", machine)
		fmt.Fprintf(w, "Name = %s\n", command)
		// get machine ip
		ip := machines[machine].Ip
		// try to connect to child
		connectChildPost(command, ip, "8080")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func connectChildPost(command string, ip string, port string) {
	fmt.Println("Trying to connect to child. " + ip + ":" + port)
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection established")
	command = command + "\n"
	fmt.Println("Sending message:\n" + command + "to child. IP: " + ip + ":" + port)
	fmt.Fprint(conn, command)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	log.Print("Server relay:", message)

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
