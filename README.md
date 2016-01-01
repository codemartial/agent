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

The Thrift example has all the auto-generated code and binaries (amd64/darwin). It can be run by itself as follows:

    [agent]$ cd thrift/gen-go/agent/server
    [server]$ go build
    [server]$ ./server &
    Starting the simple server... on  localhost:9090

Then run the benchmark client as follows:

    [server]$ cd ../client
    [client]$ go test -bench .

The Thrift C++ server binary is at `thrift/gen-cpp/servercpp`

## Results
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

On the same machine, I got the following numbers with a single non-concurrent client (note the presence of Thrift):

### Thrift Go
    $ go test -bench . -benchtime 10s -cpu 1; done
    PASS
    BenchmarkThriftClient         50000            317681 ns/op
    ok            agent/thrift/gen-go/agent/client     19.106s
### Thrift Go (Buffered Transport)
    $ go test -bench . -benchtime 10s -cpu 1; done
    PASS
    BenchmarkThriftClient	   300000	    41697 ns/op
    ok  	   agent/thrift/gen-go/agent/client	12.955s
### Thrift Go (Framed Transport)
    $ go test -bench . -benchtime 10s -cpu 1; done
    PASS
    BenchmarkThriftClient-8	  300000	     46840 ns/op
    ok  	   agent/thrift/gen-go/agent/client	14.539s    
### Thrift C++ (Buffered Transport)
    $ go test -bench . -benchtime 10s -cpu 1; done
    PASS
    BenchmarkThriftClient	   500000	    38221 ns/op
    ok  	   agent/thrift/gen-go/agent/client	19.518s
### Thrift C++ (Framed Transport)
    $ go test -bench . -benchtime 10s -cpu 1; done
    PASS
    BenchmarkThriftClient-8	  300000	     42343 ns/op
    ok  	   agent/thrift/gen-go/agent/client	13.153s
### GRPC C++
    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkGRPCClient	  200000	    113525 ns/op
    ok  	agent	23.940s
### GRPC Go
    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkGRPCClient	  200000	    105424 ns/op
    ok  	agent	22.226s
### GRPC Java
    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkGRPCClient	  100000	    108150 ns/op
    ok  	agent	12.098s

Surprisingly, `grpc-go` is faster than `grpc-java` and 2x faster than
`grpc-c++` in the concurrent scenario. In the sequential scenario,
`grpc-go` again wins, albeit with an insignificant margin. Thrift
happens to be a multi-faceted beast.

Working with Thrift is disappointing in many ways. It generates a ton
of code – 999 sloc of Go, compared to 518 of C++ with `grpc-c++` and
117 of Go with `grpc-go`. It is very difficult to figure out the
boilerplate needed to get it going (see `thrift/gen-go/agent/client`
and `thrift/gen-go/agent/server`). The client that it generates is not
concurrent, which in turn means
[_Thrift connections are not multiplexed_](https://mail-archives.apache.org/mod_mbox/thrift-user/201208.mbox/%3CA0F963DCF29346458CDF2969683DF6CC70F90B3A@SC-MBX01-2.TheFacebook.com%3E),
unlike GRPC.

The C++ backend that Thrift built out-of-the-box is single threaded. I
did find references to a non-blocking server implementation using
framed transport, but couldn't get it set up easily.

Not tabulated above, but I ran multiple client processes (instead of
invoking the same client in parallel), and got an equivalent
throughput of ~ 25000 ns/op with Thrift-Go, which is slightly better
than GRPC-Go with parallel clients. I tried emulating the test with
GRPC-Go (single client) and got the equivalent of ~ 30000 ns/op.

## Versions

    $ go version
    go version go1.5.2 darwin/amd64

    $ java -version
    java version "1.8.0_60"
    Java(TM) SE Runtime Environment (build 1.8.0_60-b27)
    Java HotSpot(TM) 64-Bit Server VM (build 25.60-b23, mixed mode)

    $ thrift -version
    Thrift version 0.9.3

    $ protoc --version
    libprotoc 3.0.0
