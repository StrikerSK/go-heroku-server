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
func (r UserRepository) ReadUserByID(userID string) (userDomains.User, bool, error) {
	return r.getUser("user_id = ?", userID)
}

// ReadUserByUsername - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUserByUsername(username string) (userDomains.User, bool, error) {
	return r.getUser("username = ?", username)
}

func (r UserRepository) getUser(query, searchValue string) (userDomains.User, bool, error) {
	var user userDomains.User
	result := r.db.Preload("Address").Where(query, searchValue).First(&user)

	if result.Error == nil {
		return user, true, nil
	}

	if result.Error == gorm.ErrRecordNotFound {
		return userDomains.User{}, false, nil
	}

	return userDomains.User{}, false, result.Error
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
