package main

import (
	"encoding/json"
	"io/ioutil"
)

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
