package board

import (
	"encoding/json"
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"
	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService ....
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) List(r ListRequest) (output *model.JSONResponse) {
	var (
		outputJSON []byte
		errorJSON  []byte
		startTime  = time.Now()
	)
	requestJSON, _ := json.Marshal(r)

	output = s.Service.List(r)
	outputJSON, _ = json.Marshal(output.Result.Data)
	errorJSON, _ = json.Marshal(output.Error)
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "List",
			"remoteAddr", r.RemoteAddr,
			"request", string(requestJSON),
			"result", string(outputJSON),
			"error", string(errorJSON),
			"took", time.Since(begin).Seconds(),
		)
	}(startTime)

	return
}

func (s *loggingService) Add(r AddRequest) (output *model.JSONResponse) {
	var (
		outputJSON []byte
		errorJSON  []byte
		startTime  = time.Now()
	)
	requestJSON, _ := json.Marshal(r)

	output = s.Service.Add(r)
	outputJSON, _ = json.Marshal(output.Result.Data)
	errorJSON, _ = json.Marshal(output.Error)
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Add",
			"remoteAddr", r.RemoteAddr,
			"request", string(requestJSON),
			"result", string(outputJSON),
			"error", string(errorJSON),
			"took", time.Since(begin).Seconds(),
		)
	}(startTime)

	return
}
