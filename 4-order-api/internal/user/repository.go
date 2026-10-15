package user

import (
	"4-order-api/pkg/db"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{Database: database}
}

func (repo *UserRepository) Create(newUser *User) (*User, error) {
	result := repo.Database.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return newUser, nil
}

func (repo *UserRepository) FindByPhone(phone string) (*User, error) {
	var searchedUser User
	result := repo.Database.DB.First(&searchedUser, User{Phone: phone})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &searchedUser, nil
}

func (repo *UserRepository) FindBySessionId(sessionId string) (*User, error) {
	var searchedUser User
	result := repo.Database.DB.First(&searchedUser, User{SessionId: sessionId})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &searchedUser, nil
}
