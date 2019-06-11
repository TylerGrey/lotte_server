package mysql

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// InitializeDatabase 데이타베이스 초기생성
func InitializeDatabase(dbName string) (*gorm.DB, error) {
	mysqlMasterConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("RDS_USERNAME"), os.Getenv("RDS_PASSWORD"), os.Getenv("RDS_HOSTNAME"), os.Getenv("RDS_PORT"), dbName)
	masterDb, err := gorm.Open("mysql", mysqlMasterConnStr)
	if err != nil {
		return nil, err
	}
	masterDb.DB()
	err = masterDb.DB().Ping()
	if err != nil {
		return nil, err
	}

	masterDb.LogMode(true)
	masterDb.DB().SetMaxIdleConns(5)
	masterDb.DB().SetMaxOpenConns(5)
	masterDb.DB().SetConnMaxLifetime(time.Minute * 5)
	masterDb.SingularTable(true)
	return masterDb, nil
}
