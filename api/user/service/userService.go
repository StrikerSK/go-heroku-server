package userServices

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/types/errors"
	userDomains "go-heroku-server/api/user/domain"
	userPorts "go-heroku-server/api/user/ports"
	"log"
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
	if _, err := s.ReadUser(user.Username); err != nil {
		if err == gorm.ErrRecordNotFound {
			user.SetRole()
			if err = s.repository.CreateUser(user); err != nil {
				log.Printf("User repository create error: %v\n", err)
				return err
			} else {
				log.Printf("User created\n")
				return nil
			}
		} else {
			log.Printf("User create error: %v\n", err)
			return err
		}
	} else {
		log.Printf("User add: user already created\n")
		return errors.NewConflictError("user already exists")
	}
}

func (s UserService) ReadUser(username string) (userDomains.User, error) {
	return s.repository.ReadUserByUsername(username)
}

func (s UserService) ReadUsers() ([]userDomains.User, error) {
	return s.repository.ReadUsers()
}

func (s UserService) UpdateUser(updatedUser userDomains.User) error {
	if err := s.repository.UpdateUser(updatedUser); err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User [%s] edit: user not found\n", updatedUser.Username)
			return err
		} else {
			log.Printf("User [%s] edit: %v\n", updatedUser.Username, err)
			return err
		}
	} else {
		return nil
	}
}
