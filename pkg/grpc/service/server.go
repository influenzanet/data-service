package service

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/influenzanet/data-service/pkg/api"
	"github.com/influenzanet/data-service/pkg/types"
	"google.golang.org/grpc"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

type dataServiceServer struct {
	clients *types.APIClients
}

// NewUserManagementServer creates a new service instance
func NewDataServiceServer(
	clients *types.APIClients,
) api.DataServiceApiServer {
	return &dataServiceServer{
		clients: clients,
	}
}

// RunServer runs gRPC service
func RunServer(ctx context.Context, port string,
	clients *types.APIClients,
) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register service
	server := grpc.NewServer()
	api.RegisterDataServiceApiServer(server, NewDataServiceServer(
		clients,
	))

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	log.Println("wait connections on port " + port)
	return server.Serve(lis)
}
