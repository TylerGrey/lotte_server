package db

import (
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/jinzhu/gorm"
)

// Board 게시글
type Board struct {
	ID        int64      `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title     string     `json:"title"`
	Image     string     `json:"image"`
	Writer    int64      `json:"writer"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// BoardRepository User 레포지터리 인터페이스
type BoardRepository interface {
	Create(board Board) model.DbChannel
	List(page int32, limit int32) model.DbChannel
}

// boardRepository 인터페이스 구조체
type boardRepository struct {
	session *gorm.DB
}

// NewBoardRepository ...
func NewBoardRepository(masterSession *gorm.DB) BoardRepository {
	return &boardRepository{
		session: masterSession,
	}
}

// Create 게시글 생성
func (r boardRepository) Create(board Board) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		err := r.session.Table("board").Create(&board).Error
		if err != nil {
			result.Err = err
		}

		result.Data = board.ID
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

// List 게시글 리스트 조회
func (r boardRepository) List(page int32, limit int32) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		boards := []*Board{}

		err := r.session.Table("board").Scan(&boards).Error
		if err != nil {
			result.Err = err
		}

		result.Data = boards
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
