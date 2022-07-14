package userServices

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	userDomains "go-heroku-server/api/user/domain"
	userPorts "go-heroku-server/api/user/ports"
)

type UserService struct {
	repository userPorts.IUserRepository
}

func NewUserService(repository userPorts.IUserRepository) UserService {
	return UserService{
		repository: repository,
	}
}

func (s UserService) CreateUser(user userDomains.User) error {
	return s.repository.CreateUser(user)
}

func (s UserService) ReadUser(username string) (userDomains.User, error) {
	return s.repository.ReadUserByUsername(username)
}

func (s UserService) ReadUsers() ([]userDomains.User, error) {
	return s.repository.ReadUsers()
}

func (s UserService) UpdateUser(updatedUser userDomains.User) error {
	return s.repository.UpdateUser(updatedUser)
}
