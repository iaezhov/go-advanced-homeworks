package auth

import (
	"4-order-api/internal/user"
)

type AuthServiceDeps struct {
	UserRepository *user.UserRepository
}
type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(deps *AuthServiceDeps) *AuthService {
	return &AuthService{
		UserRepository: deps.UserRepository,
	}
}

func (service *AuthService) Login(phone string) (sessionId string, error error) {
	existedUser, err := service.UserRepository.FindByPhone(phone)
	if err != nil {
		return "", err
	}
	if existedUser != nil {
		return existedUser.SessionId, nil
	}
	newUser := user.NewUser(phone)
	_, err = service.UserRepository.Create(newUser)
	if err != nil {
		return "", err
	}
	return newUser.SessionId, nil
}

func (service *AuthService) Verify(sessionId string, code int) (phone string, error error) {
	existedUser, err := service.UserRepository.FindBySessionId(sessionId)
	if err != nil {
		return "", err
	}
	if existedUser == nil || existedUser.Code != code {
		return "", ErrWrongCredentials
	}
	return existedUser.Phone, nil
}
