package db

import (
	"time"

	"github.com/TylerGrey/lotte_server/lib/model"

	"github.com/jinzhu/gorm"
)

// User 사용자
type User struct {
	ID        int64  `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserRepository User 레포지터리 인터페이스
type UserRepository interface {
	// 유저 등록
	Create(user User) model.DbChannel
	// 이메일로 유저 조회
	FindByEmail(email string) model.DbChannel
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
func (r userRepository) Create(user User) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		err := r.session.Table("user").Create(&user).Error
		if err != nil {
			result.Err = err
		}

		result.Data = user.ID
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

// FindByEmail 이메일로 유저 조회
func (r userRepository) FindByEmail(email string) model.DbChannel {
	storeChannel := make(model.DbChannel)
	go func() {
		result := model.DbResult{}
		user := User{}

		err := r.session.
			Table("user").
			Where("email = ?", email).
			Find(&user).Error
		if err != nil {
			result.Err = err
		}

		result.Data = user
		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}
