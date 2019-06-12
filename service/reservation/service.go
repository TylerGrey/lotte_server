package reservation

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
	ctx             context.Context
	logger          log.Logger
	reservationRepo db.ReservationRepository
}

// NewService ...
func NewService(ctx context.Context, logger log.Logger, reservationRepo db.ReservationRepository) Service {
	return &service{
		ctx:    ctx,
		logger: logger,

		reservationRepo: reservationRepo,
	}
}

// List 게시글 리스트 조회
func (s service) List(r ListRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if result := <-s.reservationRepo.List(); result.Err != nil {
		s.logger.Log("BOARD_LIST_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "게시글 조회에 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = ListResponse{
			List: result.Data.([]*model.ReservationList),
		}
	}

	return &response
}

// Add 예약 작성
func (s service) Add(r AddRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	m := db.Reservation{
		RoomID:        r.RoomID,
		BookerID:      r.UserID,
		StartDatetime: r.StartDatetime,
		EndDatetime:   r.EndDatetime,
		Title:         r.Title,
		Status:        consts.RESERVATION_WAITING,
		Attachments:   r.Attachments,
	}
	if result := <-s.reservationRepo.Create(m); result.Err != nil {
		s.logger.Log("RESERVATION_ADD_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "예약이 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = AddResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}
