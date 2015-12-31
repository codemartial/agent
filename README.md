# agent
Experiments with RPC Servers

The code here is for experimenting with very basic gRPC/Thrift stuff, except that I wrote a minimal service from scratch instead of simply running the helloworld examples bundled in the distributions.

## Using This Code

Included here are:
* a gRPC service proto, `proto/agent.proto`,
* a Go client test and benchmark `client.go`, `client_test.go`
* a Go server implementation `server/server.go`
* a C++ server implementation `server/server.cc`
* a Java server implementation `server/AgentIOServer.java` (contributed by @vidubey)

The following instructions assume that you have all the required gRPC components installed in the right place to be able to build new services. If you aren't able to build `github.com/grpc/grpc-go/examples/helloworld` and `github.com/grpc/grpc/examples/cpp/helloworld` from the [grpc-go](https://github.com/grpc/grpc-go) and [grpc](https://github.com/grpc/grpc) repos respectively, you'll most likely not be able to build the following either.
 
To build the Go stubs and skeletons, execute the following:

    [agent]$ protoc -I proto --go_out=plugins=grpc:. proto/agent.proto

Then, build the Go server:

    [agent]$ cd server
    [server]$ go build server.go

Run the server:

    [server]$ ./server -addr :9001 -nettype tcp
  
Then run the client test and benchmark:

    [server]$ cd ..
    [agent]$ go test -bench . 

To build the C++ stubs, skeletons and server, execute the following:

    [agent]$ cd server
    [server]$ make

Then run the C++ server:

    [server]$ ./servercpp
    Server listening on 0.0.0.0:9001

And re-run the benchmark as above.

The Thrift example has all the auto-generated code and can be run by itself as follows:
    [agent]$ cd thrift/gen-go/agent/server
    [server]$ go build
    [server]$ ./server &
    Starting the simple server... on  localhost:9090

Then run the benchmark client as follows:
    [server]$ cd ../client
    [client]$ go test -bench .

## Current Results
I've got the following benchmarking results on a Darwin 15.2.0, 2.5 GHz Intel Core i7 (4x2 cores w/ HT), 1600 MHz DDR3 RAM:

### C++ Server (128 Clients)
    $ go test -bench . -benchtime 10s
    PASS
    BenchmarkGRPCClient-8	  300000	     58358 ns/op
    ok  	agent	18.114s
### Go Server (128 Clients)
    $ go test -bench . -benchtime 10s
    PASS
    BenchmarkGRPCClient-8	  500000	     27313 ns/op
    ok  	agent	13.974s
### Java Server (128 Clients)
    $ go test -bench . -benchtime 10s; done
    PASS
    BenchmarkGRPCClient-8	  500000	     29197 ns/op
    ok	    agent	15.005s

On the same machine, I got the following numbers with one sequential client (note the presence of Thrift):

### Thrift (Sequential Client)
    $ go test -bench . -benchtime 10s -cpu 1; done
    PASS
    BenchmarkThriftClient	   50000	    317681 ns/op
    ok  	   agent/thrift/gen-go/agent/client	19.106s
### C++ Server (Sequential Client)
    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkGRPCClient	  200000	    113525 ns/op
    ok  	agent	23.940s
### Go Server (Sequential Client)
    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkGRPCClient	  200000	    105424 ns/op
    ok  	agent	22.226s
### Java Server (Sequential Client)
    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkGRPCClient	  100000	    108150 ns/op
    ok  	agent	12.098s

Surprisingly, `grpc-go` is faster than `grpc-java` and 2x faster than `grpc-c++` in the concurrent scenario. In the sequential scenario, `grpc-go` again wins, albeit with an insignificant margin. Thrift is just dead in the water.

Thrift is disappointing in other ways too. It generates a ton of code – 999 sloc of Go, compared to 518 of C++ with `grpc-c++` and 117 of Go with `grpc-go`. It is very difficult to figure out the boilerplate needed to get it going (see `thrift/gen-go/agent/client` and `thrift/gen-go/agent/server`. The client that it generates is not concurrent, which in turn means _Thrift_connections_are_not_multiplexed_, unlike GRPC. BTW, I did try to run 128 concurrents for Thrift, but the speedup was only about 1.2x instead of ~ 4x that is seen with `grpc-go` or `grpc-java` so I'm not sure if the Thrift server too is sufficiently concurrent.
