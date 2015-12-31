package main

import (
	"agent/thrift/gen-go/agent"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"testing"
)

var request *agent.Request
var response *agent.Response

func init() {
	setEnv()
}

func TestClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	defer client.Transport.Close()
	response, err := client.SendRequest(&agent.Request{Path: "/foo"})
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 0 || response.Status != "OK" || string(response.Body) != "/foo" {
		t.Fatal("Unexpected response:", response)
	}
}

func BenchmarkThriftClient(b *testing.B) {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Transport.Close()
	request = &agent.Request{Path: "/foo"}
	for i := 0; i < b.N; i++ {
		var err error
		response, err = client.SendRequest(request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func NewClient() (*agent.AgentIOClient, error) {
	transportFactory := env.transportFactory
	protocolFactory := env.protocolFactory
	addr := env.addr
	secure := env.secure

	var transport thrift.TTransport
	var err error
	if secure {
		return nil, fmt.Errorf("Secure transport not implemented")
	}

	transport, err = thrift.NewTSocket(addr)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return nil, err
	}
	if transport == nil {
		return nil, fmt.Errorf("Error opening socket, got nil transport. Is server available?")
	}
	transport = transportFactory.GetTransport(transport)
	if transport == nil {
		return nil, fmt.Errorf("Error from transportFactory.GetTransport(), got nil transport. Is server available?")
	}

	err = transport.Open()
	if err != nil {
		return nil, err
	}
	return agent.NewAgentIOClientFactory(transport, protocolFactory), nil
}
