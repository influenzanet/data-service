package main

import (
	"context"
	"log"

	"github.com/influenzanet/data-service/internal/config"
	gc "github.com/influenzanet/data-service/pkg/grpc/clients"
	"github.com/influenzanet/data-service/pkg/grpc/service"
	"github.com/influenzanet/data-service/pkg/types"
)

func main() {
	clients := &types.APIClients{}
	conf := config.InitConfig()

	loggingClient, close := gc.ConnectToLoggingService(conf.ServiceURLs.LoggingService)
	defer close()
	clients.LoggingService = loggingClient

	studyClient, close := gc.ConnectToStudyService(conf.ServiceURLs.StudyService)
	defer close()
	clients.StudyService = studyClient

	ctx := context.Background()
	if err := service.RunServer(
		ctx,
		conf.Port,
		clients,
	); err != nil {
		log.Fatal(err)
	}
}
