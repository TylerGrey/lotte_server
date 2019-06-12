package room

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

	Name            string `json:"name"`
	IsEnableHoliday string `json:"is_enable_holiday"`
	MinEnableTime   string `json:"min_enable_time"`
	MaxEnableTime   string `json:"max_enable_time"`
}

// UpdateRequest ...
type UpdateRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	ID              int64  `json:"id"`
	Name            string `json:"name"`
	IsEnableHoliday string `json:"is_enable_holiday"`
	MinEnableTime   string `json:"min_enable_time"`
	MaxEnableTime   string `json:"max_enable_time"`
}

// DeleteRequest ...
type DeleteRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	ID string `json:"id"`
}

// ListResponse ...
type ListResponse struct {
	List []*db.Room `json:"list"`
}

// AddResponse ...
type AddResponse struct {
	ID int64 `json:"id"`
}

// UpdateResponse ...
type UpdateResponse struct {
	ID int64 `json:"id"`
}

// DeleteResponse ...
type DeleteResponse struct {
	ID int64 `json:"id"`
}

// Endpoints ...
type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	AddEndpoint    endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
	DeleteEndpoint endpoint.Endpoint
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

// MakeUpdateEndpoint ...
func MakeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		req.UserID = int64(ctx.Value("userId").(float64))

		response := s.Update(req)
		return response, nil
	}
}

// MakeDeleteEndpoint ...
func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)

		req.UserID = int64(ctx.Value("userId").(float64))

		response := s.Delete(req)
		return response, nil
	}
}
