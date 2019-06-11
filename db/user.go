package db

import (
	"net/http"
	"time"

	"github.com/TylerGrey/lotte_server/lib/consts"
	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/TylerGrey/lotte_server/util"
	"github.com/jinzhu/gorm"
)

// User 사용자
type User struct {
	ID        int64  `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserRepository User 레포지터리 인터페이스
type UserRepository interface {
	// 유저등록
	Create(user User) chan model.DbResult
}

// userRepository 인터페이스 구조체
type userRepository struct {
	session *gorm.DB
}

// NewUserRepository ...
func NewUserRepository(masterSession *gorm.DB) UserRepository {
	return &userRepository{
		session: masterSession,
	}
}

// Create 유저 생성
func (r userRepository) Create(user User) chan model.DbResult {
	storeChannel := make(chan model.DbResult)
	go func() {
		result := model.DbResult{}
		err := r.session.Table("user").Create(&user).Error
		if err != nil {
			result.Err = util.MakeError(consts.ErrorDatabaseCode, err.Error(), http.StatusInternalServerError)
		}

		result.Data = user.ID
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
