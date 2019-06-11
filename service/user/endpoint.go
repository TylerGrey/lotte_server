package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// CreateRequest ...
type CreateRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// CreateResponse ...
type CreateResponse struct {
	ID int64 `json:"id"`
}

// Endpoints ...
type Endpoints struct {
	CreateEndpoint endpoint.Endpoint
}

// MakeCreateEndpoint ...
func MakeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)

		response := s.Create(req)
		return response, nil
	}
}
