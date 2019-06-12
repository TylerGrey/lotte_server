package db

import (
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/jinzhu/gorm"
)

// Room 회의실
type Room struct {
	ID            int64      `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Name          string     `json:"name"`
	IsEnableHoliday string     `json:"is_enable_holiday"`
	MinEnableTime string     `json:"min_enable_time"`
	MaxEnableTime string     `json:"max_enable_time"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// RoomRepository User 레포지터리 인터페이스
type RoomRepository interface {
	Create(room Room) model.DbChannel
	List(page int32, limit int32) model.DbChannel
}

// roomRepository 인터페이스 구조체
type roomRepository struct {
	session *gorm.DB
}

// NewRoomRepository ...
func NewRoomRepository(masterSession *gorm.DB) RoomRepository {
	return &roomRepository{
		session: masterSession,
	}
}

// Create 게시글 생성
func (r roomRepository) Create(room Room) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		err := r.session.Table("room").Create(&room).Error
		if err != nil {
			result.Err = err
		}

		result.Data = room.ID
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

// List 게시글 리스트 조회
func (r roomRepository) List(page int32, limit int32) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		rooms := []*Room{}

		err := r.session.Table("room").Scan(&rooms).Error
		if err != nil {
			result.Err = err
		}

		result.Data = rooms
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
