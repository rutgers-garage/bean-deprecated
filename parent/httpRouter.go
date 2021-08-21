package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
