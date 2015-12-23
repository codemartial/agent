#include <iostream>
#include <memory>
#include <string>

#include <grpc++/grpc++.h>

#include "agent.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using agent::Request;
using agent::Response;
using agent::AgentIO;

class AgentIOServiceImpl final: public AgentIO::Service {
  Status SendRequest(ServerContext* context, const Request* request,
		  Response* response) {
    response->set_status("OK");
    response->set_status_code(0);
    response->set_body(request->path());
    return Status::OK;
  }
};

void RunServer() {
  std::string server_address("0.0.0.0:9001");
  AgentIOServiceImpl service;

  ServerBuilder builder;
  // Listen on the given address without any authentication mechanism.
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  // Register "service" as the instance through which we'll communicate with
  // clients. In this case it corresponds to an *synchronous* service.
  builder.RegisterService(&service);
  // Finally assemble the server.
  std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;

  // Wait for the server to shutdown. Note that some other thread must be
  // responsible for shutting down the server for this call to ever return.
  server->Wait();
}

int main(int argc, char** argv) {
  RunServer();

  return 0;
}
