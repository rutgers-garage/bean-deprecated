package main

import (
	"fmt"
	"log"
	"net/rpc"
)

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
