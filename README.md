# agent
Experiments with gRPC

The code here is for experimenting with very basic gRPC stuff, except that I wrote a minimal service from scratch instead of simply running the helloworld examples bundled in grpc distribution.

## Using This Code

Included here are:
* a gRPC service proto, `proto/agent.proto`,
* a Go client test and benchmark `client.go`, `client_test.go`
* a Go server implementation `server/server.go`
* a C++ server implementation `server/server.cc`

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

## Current Results
I've got the following benchmarking results on a Darwin 15.2.0, 2.5 GHz Intel Core i7 (4x2 cores w/ HT), 1600 MHz DDR3 RAM:

### C++ Server
    $ go test -bench . -benchtime 10s
    PASS
    BenchmarkGRPCClient-8	  300000	     58358 ns/op
    ok  	agent	18.114s
### Go Server
    $ go test -bench . -benchtime 10s
    PASS
    BenchmarkGRPCClient-8	  500000	     27313 ns/op
    ok  	agent	13.974s

Surprisingly, grpc-go is 2x faster than grpc-c++
