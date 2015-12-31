package main

import (
	"agent/thrift/gen-go/agent"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type AgentIOImpl struct{}

func (a *AgentIOImpl) SendRequest(req *agent.Request) (r *agent.Response, err error) {
	response := agent.NewResponse()
	response.StatusCode = 0
	response.Status = "OK"
	response.Body = agent.Bytes(req.Path)
	return response, nil
}

func runServer(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TServerTransport
	var err error
	if secure {
		return fmt.Errorf("Secure transport not implemented")
	}
	transport, err = thrift.NewTServerSocket(addr)
	if err != nil {
		return err
	}
	fmt.Printf("%T\n", transport)
	handler := &AgentIOImpl{}
	processor := agent.NewAgentIOProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	fmt.Println("Starting the simple server... on ", addr)
	return server.Serve()
}
