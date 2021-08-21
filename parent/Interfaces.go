package main

type HTTPCommandBody struct {
	MachineName string
	ServiceName string
	ServiceType string
}

type BeanService struct {
	Title       string
	Endpoint    string
	ServiceType bool
}

type BeanNode struct {
	Name      string
	IP        string
	Supported []BeanService
}
