package user

import (
	"context"
	"net/http"

	"github.com/TylerGrey/lotte_server/lib/redis"

	"github.com/TylerGrey/lotte_server/lib/jwt"

	"github.com/TylerGrey/lotte_server/db"
	"github.com/TylerGrey/lotte_server/lib/consts"
	"github.com/TylerGrey/lotte_server/lib/model"
	"github.com/TylerGrey/lotte_server/util"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
)

// Service ...
type Service interface {
	SignUp(r SignUpRequest) *model.JSONResponse
	SignIn(r SignInRequest) *model.JSONResponse
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

	if util.IsEmpty(r.Email) ||
		util.IsEmpty(r.Password) ||
		util.IsEmpty(r.FirstName) ||
		util.IsEmpty(r.LastName) {
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "입력 정보를 확인해주세요.", http.StatusBadRequest)
		return &response
	}

	pw := util.GenerateFromPassword(r.Password)

	m := db.User{
		Email:     r.Email,
		Password:  pw,
		FirstName: r.FirstName,
		LastName:  r.LastName,
	}
	if result := <-s.userRepo.Create(m); result.Err != nil {
		s.logger.Log("UER_CREATE_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "회원 가입에 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = SignUpResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}

// SignIn 로그인
func (s service) SignIn(r SignInRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	// 사용자 정보 조회
	result := <-s.userRepo.FindByEmail(r.Email)
	if result.Err != nil {
		s.logger.Log("UER_FIND_BY_EMAIL_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorFindUserCode, "로그인이 실패했습니다.", http.StatusInternalServerError)
	}
	user := result.Data.(db.User)

	// 비밀번호 확인
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)); err != nil {
		s.logger.Log("COMPARE_PASSWORD_ERROR", err.Error())
		response.Error = util.MakeError(consts.ErrorValidatePasswordCode, "아이디와 비밀번호를 확인해주세요.", http.StatusInternalServerError)
	}

	// JWT 토큰 Redis에 저장
	jwtParameter := &model.Jwt{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}
	const oneWeekTime = 3600 * 24 * 7
	token, _ := jwt.GenerateToken(jwtParameter, oneWeekTime)

	redisCli := redis.MasterConnect()
	defer redisCli.Close()
	redisCli.SetValue(r.Email, token, oneWeekTime)

	response.Result.Data = SignInResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Token:     token,
	}

	return &response
}
