package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// SignUpRequest ...
type SignUpRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// SignInRequest ...
type SignInRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResponse ...
type SignUpResponse struct {
	ID int64 `json:"id"`
}

// SignInResponse ...
type SignInResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Endpoints ...
type Endpoints struct {
	SignUpEndpoint endpoint.Endpoint
	SignInEndpoint endpoint.Endpoint
}

// MakeSignUpEndpoint ...
func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)

		response := s.SignUp(req)
		return response, nil
	}
}

// MakeSignInEndpoint ...
func MakeSignInEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignInRequest)

		response := s.SignIn(req)
		return response, nil
	}
}
