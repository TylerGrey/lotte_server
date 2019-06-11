package user

import (
	"context"
	"net/http"

	"github.com/TylerGrey/lotte_server/db"
	"github.com/TylerGrey/lotte_server/lib/consts"
	"github.com/TylerGrey/lotte_server/lib/model"
	"github.com/TylerGrey/lotte_server/util"
	"github.com/go-kit/kit/log"
)

// Service ...
type Service interface {
	SignUp(r SignUpRequest) *model.JSONResponse
}

type service struct {
	ctx      context.Context
	logger   log.Logger
	userRepo db.UserRepository
}

// NewService ...
func NewService(ctx context.Context, logger log.Logger, userRepo db.UserRepository) Service {
	return &service{
		ctx:    ctx,
		logger: logger,

		userRepo: userRepo,
	}
}

// SignUp 유저 회원가입
func (s service) SignUp(r SignUpRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if (r.Email == nil || util.IsEmpty(*r.Email)) ||
		(r.Password == nil || util.IsEmpty(*r.Password)) ||
		(r.FirstName == nil || util.IsEmpty(*r.FirstName)) ||
		(r.LastName == nil || util.IsEmpty(*r.LastName)) {
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "입력 정보를 확인해주세요.", http.StatusBadRequest)
		return &response
	}

	pw := util.GenerateFromPassword(*r.Password)

	m := db.User{
		Email:     *r.Email,
		Password:  pw,
		FirstName: *r.FirstName,
		LastName:  *r.LastName,
	}
	if result := <-s.userRepo.Create(m); result.Err != nil {
		response.Error = result.Err
	} else {
		response.Result.Data = result.Data.(int64)
	}

	return &response
}
