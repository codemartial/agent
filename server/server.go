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

const gibberish1kB = `z/rt/gcAAAEDAAAAAgAAAA4AAAAICgAAAQAAAAAAAAAZAAAASAAAAF9fUEFHRVpFUk8AAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAZAAAAKAIAAF9fVEVYVAAAAAAAAAAAAAAAEAAAAAAAAAAARQAAAAAAAAAAAAAAAAAAAEUAAAAAAAcAAAAFAAAABgAAAAAAAABfX3RleHQAAAAAAAAAAAAAX19URVhUAAAAAAAAAAAAAAAgAAAAAAAA34khAAAAAAAAEAAABAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAF9fcm9kYXRhAAAAAAAAAABfX1RFWFQAAAAAAAAAAAAA4KkhAAAAAADOkBYAAAAAAOCZIQAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAX190eXBlbGluawAAAAAAAF9fVEVYVAAAAAAAAAAAAACwOjgAAAAAANhbAAAAAAAAsCo4AAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABfX2dvc3ltdGFiAAAAAAAAX19URVhUAAAAAAAAAAAAAIiWOAAAAAAAAAAAAAAAAACIhjgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAF9fZ29wY2xudGFiAAAAAABfX1RFWFQAAAAAAAAAAAAAoJY4AAAAAADTbgwAAAAAAKCGOAAFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAX19zeW1ib2xfc3R1YjEAAF9fVEVYVAAAAAAAAAAAAACABUUAAAAAANgAAAAAAAAAgPVEAAUAAAAAAAAAAAAAAAgEAIAAAAAABgAAAAAAAAAZAAAA2AEAAF9fREFUQQAAAAAAAAAAAAAAEEUAAAAAAKBgBAAAAAAAAABFAAAAAADAsAEAAAAAAAMAAAADAAAABQAAAAAAAABfX25sX3N5bWJvbF9wdHIAX19EQVRBAAAAAAAAAAAAAAAQRQAAAAAA`

func (s *AgentIOServerImpl) SendRequest(ctx context.Context, r *agent.Request) (*agent.Response, error) {
	if r.ServiceId == "large" {
		return &agent.Response{
			StatusCode: 0,
			Status:     "OK",
			Body:       []byte(gibberish1kB),
		}, nil
	}
	return &agent.Response{
		StatusCode: 0,
		Status:     "OK",
		Body:       []byte(r.Path),
	}, nil
}

func main() {
	//defer profile.Start(profile.BlockProfile).Stop()
	network := flag.String("nettype", "tcp", "network type (unix|tcp)")
	addr := flag.String("addr", ":9001", "agent server address or socket file")
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
