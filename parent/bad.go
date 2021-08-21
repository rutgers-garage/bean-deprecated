package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
)

type HTTPCommandReqBody struct {
	MachineName string
	ServiceName string
	ServiceType string
}

type HTTPBeanNodeWrapper struct {
	Name      string
	Supported []BeanService
}

type HTTPBeanServiceWrapper struct {
	Title       string
	Description string
	ServiceType string
	Params      []Param
}

type BeanNode struct {
	Name      string
	IP        string
	Supported []BeanService
}

type BeanService struct {
	Title       string
	Description string
	Endpoint    string
	ServiceType string
	Params      []Param
}

type Param struct {
	Title string
	Type  string
}

func _enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func _getBeanServiceHTTPWrapper(services []BeanService) []HTTPBeanServiceWrapper {
	var servicesWithoutEndpoint []HTTPBeanServiceWrapper

	for _, x := range services {
		servicesWithoutEndpoint = append(
			servicesWithoutEndpoint,
			HTTPBeanServiceWrapper{
				Title:       x.Title,
				Description: x.Description,
				ServiceType: x.ServiceType,
				Params:      x.Params,
			})
	}
	return servicesWithoutEndpoint
}

func _getBeanNodeHTTPWrapper(nodes []BeanNode) []HTTPBeanNodeWrapper {
	var nodesWithoutEndpoint []HTTPBeanNodeWrapper

	for _, x := range nodes {
		nodesWithoutEndpoint = append(
			nodesWithoutEndpoint,
			HTTPBeanNodeWrapper{
				Name:      x.Name,
				Supported: x.Supported,
			})
	}
	return nodesWithoutEndpoint
}

func httpGetServicesForMachine(w http.ResponseWriter, req *http.Request) {
	_enableCors(&w)
	fmt.Println("Received Web Service request")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		decoder := json.NewDecoder(req.Body)
		c := HTTPCommandReqBody{}
		err := decoder.Decode(&c)

		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(
			map[string][]HTTPBeanServiceWrapper{"services": _getBeanServiceHTTPWrapper(getServicesForMachine(c.MachineName))})

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func httpExecWebCommand(w http.ResponseWriter, req *http.Request) {
	_enableCors(&w)
	fmt.Println("Received Web Service request")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		decoder := json.NewDecoder(req.Body)
		c := HTTPCommandReqBody{}
		err := decoder.Decode(&c)

		if err != nil {
			log.Fatal(err)
		}

		machine := getMachines()[c.MachineName]
		var service string

		for _, x := range machine.Supported {
			if x.Title == c.ServiceName {
				service = x.Endpoint
			}
		}

		execWebCommand(machine.IP, service)
		fmt.Fprintf(w, "Successfully ran command")
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func httpPollMachines(w http.ResponseWriter, req *http.Request) {
	_enableCors(&w)
	fmt.Println("Received start-service request")
	switch req.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pollMachines())

	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
	}
}

func createClient(ip string) *rpc.Client {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%v:%v", ip, "8080"))
	if err != nil {
		log.Fatal("dialing:", err)
	}

	fmt.Printf("Client connected on IP %v\n", ip)
	return client
}

func pollMachines() map[string]bool {

	machines := getMachines()
	machineUpStatus := make(map[string]bool)

	for k, v := range machines {
		client := createClient(v.IP)
		childUpStatus := false
		BeanServiceErr := client.Call("Child.Up", "", &childUpStatus)

		if BeanServiceErr != nil {
			log.Print(BeanServiceErr, "(", k, ")")
		}

		machineUpStatus[v.Name] = childUpStatus
		fmt.Printf("%v response is %v\n", k, childUpStatus)
	}

	return machineUpStatus

}

func execWebCommand(ip string, endpoint string) string {
	client := createClient(ip)
	childRetString := ""

	BeanServiceErr := client.Call("Child.ExecWebService", endpoint, &childRetString)
	if BeanServiceErr != nil {
		log.Fatal(BeanServiceErr)
	}

	fmt.Printf("R response is %v\n", childRetString)

	return childRetString
}

func deserializeMachines() map[string]BeanNode {
	m := make(map[string]BeanNode)

	readBytes, _ := ioutil.ReadFile("machines.json")
	json.Unmarshal(readBytes, &m)
	print("Deserializing Machines:", m)

	return m
}

func deserializeWhitelist() map[string][]BeanService {
	m := make(map[string][]BeanService)

	readBytes, _ := ioutil.ReadFile("whitelist.json")
	json.Unmarshal(readBytes, &m)
	print("Deserializing Whitelist", m)

	return m
}

func getMachines() map[string]BeanNode {
	m := deserializeMachines()
	return m
}

func getWhitelist() map[string][]BeanService {
	m := deserializeWhitelist()
	return m
}

func getServicesForMachine(machineName string) []BeanService {
	m := deserializeMachines()
	return m[machineName].Supported
}

func main() {
	http.HandleFunc("/poll", httpPollMachines)
	http.HandleFunc("/web", httpExecWebCommand)
	http.HandleFunc("/services", httpGetServicesForMachine)

	http.ListenAndServe(":8080", nil)

}
