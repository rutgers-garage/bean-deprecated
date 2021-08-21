package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HTTPCommandReqBody struct {
	MachineName string
	ServiceName string
	ServiceType string
}

func httpExecWebCommand(w http.ResponseWriter, req *http.Request) {
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

		for i, x: range machine.Supported {
			if x.Name = c.ServiceName {
				service : x.Endpoint 
			}
		}

		execWebService(machine.IP, service)
		fmt.Fprintf(w, "Successfully ran command")
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
