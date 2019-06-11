package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/TylerGrey/lotte_server/model"
	"github.com/TylerGrey/lotte_server/util"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func toHTTPContext() httptransport.RequestFunc {
	// TODO: 추후 Auth 토큰 추가할 때 사용
	return func(ctx context.Context, r *http.Request) context.Context {
		token := r.Header.Get("Authorization")
		return context.WithValue(ctx, "Authorization", token)
	}
}

// MakeHTTPHandler ...
func MakeHTTPHandler(endpoints Endpoints, logger kitlog.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerBefore(toHTTPContext()),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}

	signUpHandler := httptransport.NewServer(
		endpoints.SignUpEndpoint,
		decodeSignUpRequest,
		encodeResponse,
		opts...,
	)

	m := mux.NewRouter()
	m.Handle("/api/user/signUp", signUpHandler).Methods("POST")

	return m
}

func decodeSignUpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := SignUpRequest{}
	request.RemoteAddr = r.RemoteAddr
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Accept", "application/json")
	w.Header().Add("server", "lotte-server")
	res := &model.JSONResponse{
		Error: &model.AppError{
			CreatedAt:  util.LocalTimeUnix(),
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		},
	}
	json.NewEncoder(w).Encode(res)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(*model.JSONResponse)
	r.Timestamp = util.LocalTimeUnix()
	if r.Error != nil {
		r.Success = false
	} else {
		r.Success = true
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
