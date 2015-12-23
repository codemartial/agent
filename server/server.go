package main

import (
	"agent"
	"flag"
	//"github.com/pkg/profile"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

type AgentIOServerImpl struct {
}

func (s *AgentIOServerImpl) SendRequest(ctx context.Context, r *agent.Request) (*agent.Response, error) {
	return &agent.Response{
		StatusCode: 0,
		Status:     "OK",
		Body:       []byte(r.Path),
	}, nil
}

func main() {
	//defer profile.Start(profile.BlockProfile).Stop()
	network := flag.String("nettype", "unix", "network type (unix|tcp)")
	addr := flag.String("addr", "/tmp/agent.sock", "agent server address or socket file")
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe(":9002", nil)) // for debug/pprof endpoint
	}()

	lis, err := net.Listen(*network, *addr)
	if err != nil {
		log.Println(err)
		return
	}

	server := grpc.NewServer()
	agent.RegisterAgentIOServer(server, &AgentIOServerImpl{})
	server.Serve(lis)
}
