package room

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
	List(r ListRequest) *model.JSONResponse
	Add(r AddRequest) *model.JSONResponse
	Update(r UpdateRequest) *model.JSONResponse
	Delete(r DeleteRequest) *model.JSONResponse
}

type service struct {
	ctx      context.Context
	logger   log.Logger
	roomRepo db.RoomRepository
}

// NewService ...
func NewService(ctx context.Context, logger log.Logger, roomRepo db.RoomRepository) Service {
	return &service{
		ctx:    ctx,
		logger: logger,

		roomRepo: roomRepo,
	}
}

// List 게시글 리스트 조회
func (s service) List(r ListRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if result := <-s.roomRepo.List(r.Page, r.Limit); result.Err != nil {
		s.logger.Log("BOARD_LIST_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "게시글 조회에 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = ListResponse{
			List: result.Data.([]*db.Room),
		}
	}

	return &response
}

// Add 게시글 작성
func (s service) Add(r AddRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	m := db.Room{}
	if result := <-s.roomRepo.Create(m); result.Err != nil {
		s.logger.Log("BOARD_ADD_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "게시글 작성 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = AddResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}

// Update 게시글 작성
func (s service) Update(r UpdateRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	m := db.Room{}
	if result := <-s.roomRepo.Create(m); result.Err != nil {
		s.logger.Log("BOARD_ADD_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "게시글 작성 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = UpdateResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}

// Delete 게시글 작성
func (s service) Delete(r DeleteRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	m := db.Room{}
	if result := <-s.roomRepo.Create(m); result.Err != nil {
		s.logger.Log("BOARD_ADD_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "게시글 작성 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = DeleteResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}
