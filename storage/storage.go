package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Open(host string, dbname string, user string, password string) (store *Storage, err error) {
	string_param := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, dbname, password)
	db, err := gorm.Open("postgres", string_param)
	store = &Storage{
		DB: db,
	}
	return
}

type Storage struct {
	DB *gorm.DB
}
