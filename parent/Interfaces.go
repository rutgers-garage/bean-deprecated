package main

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
