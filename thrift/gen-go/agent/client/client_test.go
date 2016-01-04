package main

import (
	"agent/thrift/gen-go/agent"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"sync"
	"sync/atomic"
	"testing"
)

var request *agent.Request
var response *agent.Response
var client_idx int32

const disableMT = false
const NumClients = 128
const gibberish1kB = `z/rt/gcAAAEDAAAAAgAAAA4AAAAICgAAAQAAAAAAAAAZAAAASAAAAF9fUEFHRVpFUk8AAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAAAAKAIAAF9fVEVYVAAAAAAAAAAAAAAAEAAAAAAAAAAARQAAAAAAAAAAAAAAAAAAAEUAAAAAAAcAAAAFAAAABgAAAAAAAABfX3RleHQAAAAAAAAAAAAAX19URVhUAAAAAAAAAAAAAAAgAAAAAAAA34khAAAAAAAAEAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAF9fcm9kYXRhAAAAAAAAAABfX1RFWFQAAAAAAAAAAAAA4KkhAAAAAADOkBYAAAAAAOCZIQAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAX190eXBlbGluawAAAAAAAF9fVEVYVAAAAAAAAAAAAACwOjgAAAAAANhbAAAAAAAAsCo4AAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABfX2dvc3ltdGFiAAAAAAAAX19URVhUAAAAAAAAAAAAAIiWOAAAAAAAAAAAAAAAAACIhjgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAF9fZ29wY2xudGFiAAAAAABfX1RFWFQAAAAAAAAAAAAAoJY4AAAAAADTbgwAAAAAAKCGOAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAX19zeW1ib2xfc3R1YjEAAF9fVEVYVAAAAAAAAAAAAACABUUAAAAAANgAAAAAAAAAgPVEAAUAAAAAAAAAAAAAAAgEAIAAAAAABgAAAAAAAAAZAAAA2AEAAF9fREFUQQAAAAAAAAAAAAAAEEUAAAAAAKBgBAAAAAAAAABFAAAAAADAsAEAAAAAAAMAAAADAAAABQAAAAAAAABfX25sX3N5bWJvbF9wdHIAX19EQVRBAAAAAAAAAAAAAAAQRQAAAAAA`

var clients = [NumClients]*agent.AgentIOClient{}
var mus = [NumClients]sync.Mutex{}

func init() {
	setEnv()
	if disableMT {
		return
	}

	for i := range clients {
		client, err := NewClient()
		if err != nil {
			panic(err)
		}
		clients[i] = client
	}
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
	if disableMT {
		benchmarkThriftClientST(b, false)
	} else {
		benchmarkThriftClientMT(b, false)
	}
}

func BenchmarkThriftClientLarge(b *testing.B) {
	if disableMT {
		benchmarkThriftClientST(b, true)
	} else {
		benchmarkThriftClientMT(b, true)
	}
}

func benchmarkThriftClientST(b *testing.B, large bool) {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Transport.Close()
	request = &agent.Request{Path: "/foo"}
	if large {
		request.ServiceID = "large"
		request.Body = agent.Bytes(gibberish1kB)
	}
	for i := 0; i < b.N; i++ {
		var err error
		response, err = client.SendRequest(request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkThriftClientMT(b *testing.B, large bool) {
	request = &agent.Request{Path: "/foo"}
	if large {
		request.ServiceID = "large"
		request.Body = agent.Bytes(gibberish1kB)
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idx := atomic.AddInt32(&client_idx, 1) % NumClients
			mus[idx].Lock()
			client := clients[idx]
			var err error
			response, err = client.SendRequest(request)
			if err != nil {
				mus[idx].Unlock()
				b.Fatal(err)
			}
			mus[idx].Unlock()
		}
	})
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
