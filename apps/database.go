package apps

import (
	"fmt"
	"os"

	"github.com/madjiebimaa/go-basic-auth/helpers"
	"github.com/madjiebimaa/go-basic-auth/models/domains"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	DB_NAME := os.Getenv("DB_NAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(`%s:%s@/go_basic_auth`, DB_NAME, DB_PASSWORD)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	helpers.PanicIfError(err)

	db.AutoMigrate(&domains.User{})

	return db
}
