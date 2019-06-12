package reservation

import (
	"context"
	"net/http"

	"github.com/TylerGrey/lotte_server/util"

	"github.com/TylerGrey/lotte_server/lib/jwt"
	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/go-kit/kit/endpoint"
)

// ListRequest ...
type ListRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`
}

// AddRequest ...
type AddRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	RoomID        int64   `json:"roomId"`
	Users         []int64 `json:"userIds"`
	StartDatetime string  `json:"startDatetime"`
	EndDatetime   string  `json:"endDatetime"`
	Title         string  `json:"title"`
	Attachments   string  `json:"attachments"`
}

// FindRequest ...
type FindRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	ID int64 `json:"id"`
}

// UpdateStatusRequest ...
type UpdateStatusRequest struct {
	UserID     int64  `json:"userId"`
	RemoteAddr string `json:"remoteAddr"`

	ID     int64  `json:"id"`
	Status string `json:"status"`
}

// ListResponse ...
type ListResponse struct {
	List []*model.ReservationList `json:"list"`
}

// AddResponse ...
type AddResponse struct {
	ID int64 `json:"id"`
}

// FindResponse ...
type FindResponse struct {
	Reservation *model.ReservationList `json:"reservation"`
}

// UpdateStatusResponse
type UpdateStatusResponse struct {
	ID int64 `json:"id"`
}

// Endpoints ...
type Endpoints struct {
	ListEndpoint         endpoint.Endpoint
	AddEndpoint          endpoint.Endpoint
	FindEndpoint         endpoint.Endpoint
	UpdateStatusEndpoint endpoint.Endpoint
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

// MakeFindEndpoint ...
func MakeFindEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRequest)

		response := s.Find(req)
		return response, nil
	}
}

// MakeUpdateStatusEndpoint ...
func MakeUpdateStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateStatusRequest)

		response := s.UpdateStatus(req)
		return response, nil
	}
}
