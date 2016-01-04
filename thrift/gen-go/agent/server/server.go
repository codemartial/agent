package main

import (
	"agent/thrift/gen-go/agent"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type AgentIOImpl struct{}

const gibberish1kB = `z/rt/gcAAAEDAAAAAgAAAA4AAAAICgAAAQAAAAAAAAAZAAAASAAAAF9fUEFHRVpFUk8AAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAAAAKAIAAF9fVEVYVAAAAAAAAAAAAAAAEAAAAAAAAAAARQAAAAAAAAAAAAAAAAAAAEUAAAAAAAcAAAAFAAAABgAAAAAAAABfX3RleHQAAAAAAAAAAAAAX19URVhUAAAAAAAAAAAAAAAgAAAAAAAA34khAAAAAAAAEAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAF9fcm9kYXRhAAAAAAAAAABfX1RFWFQAAAAAAAAAAAAA4KkhAAAAAADOkBYAAAAAAOCZIQAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAX190eXBlbGluawAAAAAAAF9fVEVYVAAAAAAAAAAAAACwOjgAAAAAANhbAAAAAAAAsCo4AAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABfX2dvc3ltdGFiAAAAAAAAX19URVhUAAAAAAAAAAAAAIiWOAAAAAAAAAAAAAAAAACIhjgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAF9fZ29wY2xudGFiAAAAAABfX1RFWFQAAAAAAAAAAAAAoJY4AAAAAADTbgwAAAAAAKCGOAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAX19zeW1ib2xfc3R1YjEAAF9fVEVYVAAAAAAAAAAAAACABUUAAAAAANgAAAAAAAAAgPVEAAUAAAAAAAAAAAAAAAgEAIAAAAAABgAAAAAAAAAZAAAA2AEAAF9fREFUQQAAAAAAAAAAAAAAEEUAAAAAAKBgBAAAAAAAAABFAAAAAADAsAEAAAAAAAMAAAADAAAABQAAAAAAAABfX25sX3N5bWJvbF9wdHIAX19EQVRBAAAAAAAAAAAAAAAQRQAAAAAA`

func (a *AgentIOImpl) SendRequest(req *agent.Request) (r *agent.Response, err error) {
	response := agent.NewResponse()
	response.StatusCode = 0
	response.Status = "OK"
	if req.ServiceID == "large" {
		response.Body = agent.Bytes(gibberish1kB)
	} else {
		response.Body = agent.Bytes(req.Path)
	}
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
