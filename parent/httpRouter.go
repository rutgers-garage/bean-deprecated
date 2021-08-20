package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func startService(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Received start-service request")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		decoder := json.NewDecoder(req.Body)
		c := AWSServiceRequestBody{}
		err := decoder.Decode(&c)

		if err != nil {
			log.Fatal(err)
		}

		startAWSService(c)
		fmt.Fprintf(w, "Successfully started new service")
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
