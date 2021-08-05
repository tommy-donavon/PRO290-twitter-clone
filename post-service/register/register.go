package register

import (
	"fmt"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

type (
	ConsulClient struct {
		c           *consulapi.Client
		ServiceName string
	}
	Service struct {
		Address string
		Port    int
		ID      string
	}
)

func NewConsulClient(serviceName string) *ConsulClient {
	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		panic(err)
	}
	return &ConsulClient{consul, serviceName}
}

func (client *ConsulClient) RegisterService() error {
	reg := new(consulapi.AgentServiceRegistration)
	reg.Name = client.ServiceName
	reg.ID = hostname()
	reg.Address = hostname()
	port, err := strconv.Atoi(os.Getenv("PORT")[1:len(os.Getenv("PORT"))])
	if err != nil {
		return err
	}
	reg.Port = port
	reg.Check = new(consulapi.AgentServiceCheck)
	reg.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", hostname(), port)
	reg.Check.Interval = "5s"
	reg.Check.Timeout = "3s"
	return client.c.Agent().ServiceRegister(reg)

}

func (client *ConsulClient) DeregisterService() error {
	return client.c.Agent().ServiceDeregister(hostname())
}

func (client *ConsulClient) LookUpService(serviceId string) (*Service, error) {
	services, err := client.c.Agent().Services()
	if err != nil {
		return nil, err
	}
	for _, k := range services {
		if k.Service == serviceId {
			return &Service{k.Address, k.Port, k.ID}, nil
		}
	}
	return nil, fmt.Errorf("No service found")

}

func (service *Service) GetHTTP() string {
	return fmt.Sprintf("http://%s:%v/", service.Address, service.Port)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hn
}
