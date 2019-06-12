package db

import (
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/jinzhu/gorm"
)

// Reservation 예약
type Reservation struct {
	ID            int64      `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	RoomID        int64      `json:"room_id"`
	BookerID      int64      `json:"booker_id"`
	StartDatetime string     `json:"start_datetime"`
	EndDatetime   string     `json:"end_datetime"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	Attachments   string     `json:"attachments"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// ReservationRepository User 레포지터리 인터페이스
type ReservationRepository interface {
	Create(reservation Reservation) model.DbChannel
	List() model.DbChannel
}

// reservationRepository 인터페이스 구조체
type reservationRepository struct {
	session *gorm.DB
}

// NewReservationRepository ...
func NewReservationRepository(masterSession *gorm.DB) ReservationRepository {
	return &reservationRepository{
		session: masterSession,
	}
}

// Create 예약 생성
func (r reservationRepository) Create(reservation Reservation) model.DbChannel {
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
func (r reservationRepository) List() model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		rooms := []*model.ReservationList{}

		tx := r.session.Table("reservation r, user u")
		tx = tx.Select("r.id, room_id, booker_id, start_datetime, end_datetime, title, name AS user_name, attendee_name")
		tx = tx.Where("r.booker_id = u.id")
		err := tx.Scan(&rooms).Error
		if err != nil {
			result.Err = err
		}

		result.Data = rooms
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
