package main

import "net/http"

func main() {
	http.HandleFunc("/poll", httpPollMachines)
	http.HandleFunc("/web", httpExecWebCommand)
	http.HandleFunc("/services", httpGetServicesForMachine)

	http.ListenAndServe(":8080", nil)

}
