package userRepositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-heroku-server/api/src/errors"
	userDomains "go-heroku-server/api/user/domain"
)

type UserRepository struct {
	Repository *gorm.DB
}

func NewUserRepository(repository *gorm.DB) UserRepository {
	repository.AutoMigrate(&userDomains.User{}, &userDomains.Address{})
	return UserRepository{
		Repository: repository,
	}
}

func (r UserRepository) CreateUser(user userDomains.User) (err error) {
	if err = r.Repository.Create(&user).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

// ReadUsers - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUsers() (user []userDomains.User, err error) {
	if err = r.Repository.Preload("Address").Find(&user).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

// ReadUserByID - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUserByID(userID string) (user userDomains.User, err error) {
	if err = r.Repository.Preload("Address").Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("user [%s] not found", userID))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}

// ReadUserByUsername - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUserByUsername(username string) (user userDomains.User, err error) {
	if err = r.Repository.Preload("Address").Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("user [%s] not found", username))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r UserRepository) UpdateUser(updatedUser userDomains.User) (err error) {
	if err = r.Repository.Model(&userDomains.User{}).Where("username = ?", updatedUser.Username).Update(&updatedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("user [%s] not found", updatedUser.Username))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}
