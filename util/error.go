package util

import (
	"net/http"

	"github.com/TylerGrey/lotte_server/lib/consts"

	"github.com/TylerGrey/lotte_server/model"
)

func _makeErrorCode(code int32, err model.AppError) *model.AppError {
	switch code {
	case consts.ErrorBadRequestCode:
		{
			err.ErrorCode = consts.ErrorBadRequestCode
			err.ErrorMsg = consts.ErrorBadRequestMsg
		}
	case consts.ErrorDatabaseCode:
		{
			err.ErrorCode = consts.ErrorDatabaseCode
			err.ErrorMsg = consts.ErrorDatabaseMsg
		}
	default:
	}

	return &err
}

// MakeError 에러 생성
func MakeError(code int32, message string, statusCode int32) *model.AppError {
	err := &model.AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		CreatedAt:  LocalTimeUnix(),
	}

	return _makeErrorCode(code, *err)
}
