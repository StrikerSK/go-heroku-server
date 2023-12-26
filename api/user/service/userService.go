package userServices

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/src/errors"
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
	_, exist, err := s.repository.ReadUserByUsername(user.Username)

	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}

	if !exist {
		user.SetRole()
		user.EncryptPassword()
		if err = s.repository.CreateUser(user); err != nil {
			log.Printf("User repository create error: %v\n", err)
			return err
		} else {
			log.Println("User created!")
			return nil
		}
	} else {
		log.Printf("User add: user already created\n")
		return errors.NewConflictError("user already exists")
	}
}

func (s UserService) ReadUser(username string) (userDomains.User, error) {
	user, exists, err := s.repository.ReadUserByUsername(username)

	if err != nil {
		return userDomains.User{}, errors.NewBadRequestError(err.Error())
	}

	if !exists {
		return userDomains.User{}, errors.NewNotFoundError("user does not exist")
	}

	return user, nil
}

func (s UserService) ReadUsers() ([]userDomains.User, error) {
	users, err := s.repository.ReadUsers()
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].ClearPassword()
	}

	return users, nil
}

func (s UserService) UpdateUser(updatedUser userDomains.User) error {
	updatedUser.EncryptPassword()
	if err := s.repository.UpdateUser(updatedUser); err != nil {
		_, notFoundErr := err.(errors.NotFoundError)
		if notFoundErr {
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
