package userRepositories

import (
	"fmt"
	"go-heroku-server/api/src/errors"
	userDomains "go-heroku-server/api/user/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	_ = db.AutoMigrate(&userDomains.User{}, &userDomains.Address{})
	return UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateUser(user userDomains.User) (err error) {
	if err = r.db.Create(&user).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

// ReadUsers - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUsers() (user []userDomains.User, err error) {
	if err = r.db.Preload("Address").Find(&user).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

// ReadUserByID - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUserByID(userID string) (user userDomains.User, err error) {
	if err = r.db.Preload("Address").Where("user_id = ?", userID).First(&user).Error; err != nil {
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
	if err = r.db.Preload("Address").Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("user [%s] not found", username))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r UserRepository) UpdateUser(updatedUser userDomains.User) (err error) {
	if err = r.db.Model(&userDomains.User{}).Where("username = ?", updatedUser.Username).Updates(&updatedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("user [%s] not found", updatedUser.Username))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}
