package db

import (
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/jinzhu/gorm"
)

// Attendee 예약
type Attendee struct {
	ID            int64      `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	ReservationID int64      `json:"reservation_id"`
	UserID        int64      `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// AttendeeRepository User 레포지터리 인터페이스
type AttendeeRepository interface {
	Create(reservation Attendee) model.DbChannel
	List(reservationID int32) model.DbChannel
}

// attendeeRepository 인터페이스 구조체
type attendeeRepository struct {
	session *gorm.DB
}

// NewAttendeeRepository ...
func NewAttendeeRepository(masterSession *gorm.DB) AttendeeRepository {
	return &attendeeRepository{
		session: masterSession,
	}
}

// Create 예약 생성
func (r attendeeRepository) Create(reservation Attendee) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		err := r.session.Table("reservation").Create(&reservation).Error
		if err != nil {
			result.Err = err
		}

		result.Data = reservation.ID
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

// List 예약 리스트 조회
func (r attendeeRepository) List(reservationID int32) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		rooms := []*Attendee{}

		err := r.session.Table("reservation").Where("reservation_id = ?", reservationID).Scan(&rooms).Error
		if err != nil {
			result.Err = err
		}

		result.Data = rooms
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
