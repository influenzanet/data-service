package types

import (
	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
	studyAPI "github.com/influenzanet/study-service/pkg/api"
)

// APIClients holds the service clients to the internal services
type APIClients struct {
	LoggingService loggingAPI.LoggingServiceApiClient
	StudyService   studyAPI.StudyServiceApiClient
}
