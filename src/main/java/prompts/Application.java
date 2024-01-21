package prompts;

import io.grpc.Server;
import io.grpc.protobuf.services.ProtoReflectionService;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.core.env.Environment;

@SpringBootApplication
public class Application {
    public static void main(String[] args) throws Exception {
        SpringApplication app = new SpringApplication(Application.class);
        Environment env = app.run(args).getEnvironment();
        
        int port = Integer.parseInt(env.getProperty("server.port"));
        Server server = io.grpc.ServerBuilder.forPort(port)
            .addService(new TransferServiceImpl(env))
            .addService(ProtoReflectionService.newInstance()) // Add reflection service
            .build();
        
        server.start();
        System.out.println("Server started, listening on " + port);
        server.awaitTermination();
    }
}
