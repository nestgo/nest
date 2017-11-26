package nest

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcServer *grpc.Server

//Registeration Register 方法
type Registeration interface {
	Register(*grpc.Server)
}

//Register register services
func Register(services ...interface{}) {
	grpcServer = grpc.NewServer()

	for _, service := range services {
		if registerService, _ := service.(Registeration); registerService != nil {
			registerService.Register(grpcServer)
		}
	}

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
}

//Run start grpc server
func Run(address string) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("Listener at %s", address)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
