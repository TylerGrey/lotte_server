package reservation

import (
	"context"
	"net/http"
	"strings"

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
	Find(r FindRequest) *model.JSONResponse
	UpdateStatus(r UpdateStatusRequest) *model.JSONResponse
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

// List 예약 리스트 조회
func (s service) List(r ListRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if result := <-s.reservationRepo.List(); result.Err != nil {
		s.logger.Log("BOARD_LIST_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "예약 조회에 실패했습니다.", http.StatusInternalServerError)
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

// Find 예약 리스트 조회
func (s service) Find(r FindRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if result := <-s.reservationRepo.Find(r.ID); result.Err != nil {
		s.logger.Log("BOARD_LIST_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "예약 조회에 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = FindResponse{
			Reservation: result.Data.(*model.ReservationList),
		}
	}

	return &response
}

// UpdateStatus 예약 상태 변경
func (s service) UpdateStatus(r UpdateStatusRequest) *model.JSONResponse {
	response := model.JSONResponse{}

	if strings.Compare(r.Status, "APPROVE") == 0 {
		// TODO: 예약 안겹치는 로직 추가
	} else if strings.Compare(r.Status, "REJECT") != 0 {
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "필수 항목 값이 없습니다.", http.StatusBadRequest)
		return &response
	}

	m := db.Reservation{
		ID:     r.ID,
		Status: r.Status,
	}
	if result := <-s.reservationRepo.Update(m); result.Err != nil {
		s.logger.Log("UPDATE_STATUS_ERROR", result.Err.Error())
		response.Error = util.MakeError(consts.ErrorBadRequestCode, "예약 상태 변경에 실패했습니다.", http.StatusInternalServerError)
	} else {
		response.Result.Data = UpdateStatusResponse{
			ID: result.Data.(int64),
		}
	}

	return &response
}
