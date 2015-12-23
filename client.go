package agent

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type ClientFactory func(string, string) AgentIOClient

var Client ClientFactory = func() ClientFactory {
	var agentIOClient AgentIOClient = nil
	return func(network, addr string) AgentIOClient {
		if agentIOClient != nil {
			return agentIOClient
		}
		var once sync.Once
		once.Do(func() {
			if c, err := NewClient(network, addr); err != nil {
				return
			} else {
				agentIOClient = c
			}
		})
		return agentIOClient
	}
}()

func NewClient(network, addr string) (AgentIOClient, error) {
	cc, err := grpc.Dial(addr,
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.Dial(network, addr)
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return NewAgentIOClient(cc), nil
}
