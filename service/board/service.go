package board

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
}

type service struct {
	ctx       context.Context
	logger    log.Logger
	boardRepo db.BoardRepository
}

// NewService ...
func NewService(ctx context.Context, logger log.Logger, boardRepo db.BoardRepository) Service {
	return &service{
		ctx:    ctx,
		logger: logger,

		boardRepo: boardRepo,
	}
}

// List 게시글 리스트 조회
func (s service) List(r ListRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if result := <-s.boardRepo.List(r.Page, r.Limit); result.Err != nil {
		s.logger.Log("BOARD_LIST_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "게시글 조회에 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = ListResponse{
			List: result.Data.([]*db.Board),
		}
	}

	return &response
}

// Add 게시글 작성
func (s service) Add(r AddRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	m := db.Board{
		Title:  r.Title,
		Writer: r.UserID,
	}
	if result := <-s.boardRepo.Create(m); result.Err != nil {
		s.logger.Log("BOARD_ADD_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorCreateUserCode, "게시글 작성 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = AddResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}
