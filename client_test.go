package agent_test

import (
	"agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

var client_idx int32

var network, addr string = "tcp", "127.0.0.1:9001"

//var network, addr string = "unix", "/tmp/agent.sock"

var client agent.AgentIOClient
var request *agent.Request
var response *agent.Response

const NumClients = 128

type ClientPool [NumClients]agent.AgentIOClient

var clients ClientPool = ClientPool{}

func init() {
	log.SetFlags(log.Lshortfile)
	for i := range clients {
		client, err := agent.NewClient(network, addr)
		if err != nil {
			panic(err)
		}
		clients[i] = client
	}
}

func TestClient(t *testing.T) {
	cc, err := grpc.Dial(addr,
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.Dial(network, addr)
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	defer cc.Close()

	client = agent.NewAgentIOClient(cc)
	response, err := client.SendRequest(context.Background(), &agent.Request{Path: "/foo"})
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 0 || response.Status != "OK" || string(response.Body) != "/foo" {
		t.Fatal("Unexpected response:", response.String())
	}
}

func TestOtherClient(t *testing.T) {
	client = agent.Client(network, addr)
	response, err := client.SendRequest(context.Background(), &agent.Request{Path: "/foo"})
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 0 || response.Status != "OK" || string(response.Body) != "/foo" {
		t.Fatal("Unexpected response:", response)
	}
}

func BenchmarkGRPCClient(b *testing.B) {
	runBenchmark(b, false)
}

func BenchmarkGRPCClientLarge(b *testing.B) {
	runBenchmark(b, true)
}

func runBenchmark(b *testing.B, large bool) {
	request = &agent.Request{Path: "/foo"}
	if large {
		request.ServiceId = "large"
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client := clients[atomic.LoadInt32(&client_idx)%NumClients]
			atomic.AddInt32(&client_idx, 1)
			var err error
			response, err = client.SendRequest(context.Background(), request)
			// Ignore racy write to "response".
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
