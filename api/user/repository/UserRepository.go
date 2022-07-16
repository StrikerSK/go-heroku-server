package userRepositories

import (
	"github.com/jinzhu/gorm"
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

func (r UserRepository) CreateUser(user userDomains.User) error {
	return r.Repository.Create(&user).Error
}

// ReadUsers - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUsers() (user []userDomains.User, err error) {
	err = r.Repository.Preload("Address").Find(&user).Error
	return
}

// ReadUserByID - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUserByID(userID string) (user userDomains.User, err error) {
	err = r.Repository.Preload("Address").Where("username = ?", userID).First(&user).Error
	return
}

// ReadUserByUsername - retrieves user and flag if exists can be registered to database
func (r UserRepository) ReadUserByUsername(username string) (user userDomains.User, err error) {
	err = r.Repository.Preload("Address").Where("username = ?", username).First(&user).Error
	return
}

func (r UserRepository) UpdateUser(updatedUser userDomains.User) (err error) {
	err = r.Repository.Model(&userDomains.User{}).Update(&updatedUser).Error
	return
}
