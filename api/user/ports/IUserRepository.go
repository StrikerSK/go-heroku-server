package userPorts

import (
	userDomains "go-heroku-server/api/user/domain"
)

type IUserRepository interface {
	CreateUser(userDomains.User) error
	ReadUsers() ([]userDomains.User, error)
	ReadUserByID(string) (userDomains.User, bool, error)
	ReadUserByUsername(string) (userDomains.User, bool, error)
	UpdateUser(userDomains.User) error
}
