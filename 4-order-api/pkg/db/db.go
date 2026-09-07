package db

import (
	"4-order-api/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
