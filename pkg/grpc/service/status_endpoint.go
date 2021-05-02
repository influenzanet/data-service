package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/influenzanet/go-utils/pkg/api_types"
)

// Status endpoint should return internal status of the system if running correctly
func (s *dataServiceServer) Status(ctx context.Context, _ *empty.Empty) (*api_types.ServiceStatus, error) {
	return &api_types.ServiceStatus{
		Status:  api_types.ServiceStatus_NORMAL,
		Msg:     "service running",
		Version: apiVersion,
	}, nil
}
