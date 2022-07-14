package userPorts

import (
	userDomains "go-heroku-server/api/user/domain"
)

type IUserService interface {
	CreateUser(userDomains.User) error
	ReadUser(string) (userDomains.User, error)
	ReadUsers() ([]userDomains.User, error)
	UpdateUser(userDomains.User) error
}
