package board

import (
	"context"
	"net/http"

	"github.com/TylerGrey/lotte_server/util"

	"github.com/TylerGrey/lotte_server/db"
	"github.com/TylerGrey/lotte_server/lib/jwt"
	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/go-kit/kit/endpoint"
)

// ListRequest ...
type ListRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	Page  int32 `json:"page"`
	Limit int32 `json:"total"`
}

// AddRequest ...
type AddRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	Title string `json:"title"`
}

// ListResponse ...
type ListResponse struct {
	List []*db.Board `json:"list"`
}

// AddResponse ...
type AddResponse struct {
	ID int64 `json:"id"`
}

// Endpoints ...
type Endpoints struct {
	ListEndpoint endpoint.Endpoint
	AddEndpoint  endpoint.Endpoint
}

// MakeAuthVerifyMiddleware ...
func MakeAuthVerifyMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			authorization := ctx.Value("Authorization").(string)
			ctx, err = jwt.RequireTokenAuthentication(ctx, authorization)
			if err != nil {
				return &model.JSONResponse{
					Error: util.MakeError(5000, err.Error(), http.StatusBadRequest),
				}, nil
			}

			return next(ctx, request)
		}
	}
}

// MakeListEndpoint ...
func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)

		response := s.List(req)
		return response, nil
	}
}

// MakeAddEndpoint ...
func MakeAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)

		req.UserID = int64(ctx.Value("userId").(float64))

		response := s.Add(req)
		return response, nil
	}
}