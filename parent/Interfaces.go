package main

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
