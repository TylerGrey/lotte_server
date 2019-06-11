package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// SignUpRequest ...
type SignUpRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	Email     *string `json:"email"`
	Password  *string `json:"password"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

// SignUpResponse ...
type SignUpResponse struct {
	ID int64 `json:"id"`
}

// Endpoints ...
type Endpoints struct {
	SignUpEndpoint endpoint.Endpoint
}

// MakeSignUpEndpoint ...
func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)

		response := s.SignUp(req)
		return response, nil
	}
}
