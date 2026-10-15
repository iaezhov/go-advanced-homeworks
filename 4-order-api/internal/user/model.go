package user

import (
	"math/rand"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone     string `gorm:"index"`
	SessionId string `gorm:"index"`
	Code      int
}

func NewUser(phone string) *User {
	user := &User{
		Phone: phone,
	}
	user.GenerateSessionId()
	user.GenerateCode()
	return user
}

func (u *User) GenerateSessionId() {
	u.SessionId = RandStringRunes(10)
}

func (u *User) GenerateCode() {
	u.Code = rand.Intn(9000) + 1000
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
