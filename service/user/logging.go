package user

import (
	"encoding/json"
	"time"

	"github.com/TylerGrey/lotte_server/model"
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

func (s *loggingService) SignUp(r SignUpRequest) (output *model.JSONResponse) {
	var (
		outputJSON []byte
		errorJSON  []byte
		startTime  = time.Now()
	)
	requestJSON, _ := json.Marshal(r)

	output = s.Service.SignUp(r)
	outputJSON, _ = json.Marshal(output.Result.Data)
	errorJSON, _ = json.Marshal(output.Error)
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SignUp",
			"remoteAddr", r.RemoteAddr,
			"request", string(requestJSON),
			"result", string(outputJSON),
			"error", string(errorJSON),
			"took", time.Since(begin).Seconds(),
		)
	}(startTime)

	return
}
