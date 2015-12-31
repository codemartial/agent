package agent;

import io.grpc.Server;
import io.grpc.netty.NettyServerBuilder;
import io.grpc.stub.StreamObserver;

import java.io.IOException;
import java.util.logging.Logger;

/**
 * Created by vidyakant.dubey on 24/12/15.
 */
public class AgentIOServer {
    private static final Logger logger = Logger.getLogger(AgentIOServer.class.getName());

    private final int port;
    private Server server;

    public AgentIOServer(int port){
        this.port=port;
    }

    public void start() throws IOException {
        server = NettyServerBuilder.forPort(port)
                .addService(AgentIOGrpc.bindService(new AgentIOServiceImpl()))
                .build()
                .start();
        logger.info("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                System.err.println("*** shutting down gRPC server since JVM is shutting down");
                AgentIOServer.this.stop();
                System.err.println("*** server shut down");
            }
        });
    }

    public void stop() {
        if (server != null) {
            server.shutdown();
        }
    }

    private void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) throws Exception {
        AgentIOServer server = new AgentIOServer(9001);
        server.start();
        server.blockUntilShutdown();
    }

    private class AgentIOServiceImpl implements AgentIOGrpc.AgentIO {
        @Override
        public void sendRequest(Sample.Request request, StreamObserver<Sample.Response> responseObserver) {
            Sample.Response response = Sample.Response.newBuilder()
                    .setStatusCode(0)
                    .setStatus("OK")
                    .setBody(request.getPathBytes())
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
        }
    }

}
