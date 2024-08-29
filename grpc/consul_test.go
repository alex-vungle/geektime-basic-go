package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

type ConsulRegistryTestSuite struct {
	suite.Suite
	client *consul.Client
}

func (s *ConsulRegistryTestSuite) SetupSuite() {
	client, err := consul.NewClient(&consul.Config{
		Address: "localhost:8500",
	})
	require.NoError(s.T(), err)
	s.client = client
}

func (s *ConsulRegistryTestSuite) TestServer() {
	l, err := net.Listen("tcp", ":8090")
	require.NoError(s.T(), err)
	err = s.consulRegister()
	require.NoError(s.T(), err)
	server := grpc.NewServer()
	RegisterUserServiceServer(server, &Server{})
	server.Serve(l)
}

func (s *ConsulRegistryTestSuite) consulRegister() error {
	//生成注册对象
	registration := &consul.AgentServiceRegistration{
		Name:    "user-service",
		ID:      uuid.New().String(),
		Port:    8090,
		Address: "127.0.0.1",
	}
	err := s.client.Agent().ServiceRegister(registration)
	return err
}

func (s *ConsulRegistryTestSuite) TestClient() {
	t := s.T()
	// 解析服务的地址
	address, err := resolveConsulService("user-service", "localhost:8500")
	require.NoError(t, err)

	cc, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	//cc, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	require.NoError(t, err)
	client := NewUserServiceClient(cc)
	resp, err := client.GetByID(context.Background(), &GetByIDRequest{Id: 123})
	require.NoError(t, err)
	t.Log(resp.User)
}

func resolveConsulService(serviceName string, consulAddr string) (string, error) {
	client, err := consul.NewClient(&consul.Config{Address: consulAddr})
	if err != nil {
		return "", err
	}

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	service := services[0]
	address := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	return address, nil
}

func TestConsul(t *testing.T) {
	suite.Run(t, new(ConsulRegistryTestSuite))
}
