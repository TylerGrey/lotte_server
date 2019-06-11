package user

import (
	"context"
	"net/http"

	"github.com/TylerGrey/lotte_server/lib/consts"
	"github.com/TylerGrey/lotte_server/model"
	"github.com/TylerGrey/lotte_server/model/user"
	"github.com/TylerGrey/lotte_server/util"
	"github.com/go-kit/kit/log"
)

// Service ...
type Service interface {
	Create(r CreateRequest) *model.JSONResponse
}

type service struct {
	ctx      context.Context
	logger   log.Logger
	userRepo user.Repository
}

// NewService ...
func NewService(ctx context.Context, logger log.Logger, userRepo user.Repository) Service {
	return &service{
		ctx:    ctx,
		logger: logger,

		userRepo: userRepo,
	}
}

// Create 유저 회원가입
func (s service) Create(r CreateRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if r.Email == nil || r.Password == nil {
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "입력 정보를 확인해주세요.", http.StatusBadRequest)
		return &response
	}

	pw := util.GenerateFromPassword(*r.Password)

	m := user.User{
		Email:    *r.Email,
		Password: pw,
	}
	if result := <-s.userRepo.Create(m); result.Err != nil {
		response.Error = result.Err
	} else {
		response.Result.Data = result.Data.(int64)
	}

	return &response
}
