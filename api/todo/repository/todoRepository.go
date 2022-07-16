package todoRepositories

import (
	"github.com/jinzhu/gorm"
	"go-heroku-server/api/todo/domain"
)

type TodoRepository struct {
	repository *gorm.DB
}

func NewTodoRepository(repository *gorm.DB) TodoRepository {
	repository.AutoMigrate(todoDomains.Todo{})
	return TodoRepository{
		repository: repository,
	}
}

func (r TodoRepository) CreateTodo(newTodo todoDomains.Todo) error {
	return r.repository.Create(&newTodo).Error
}

func (r TodoRepository) ReadTodo(todoID uint) (todo todoDomains.Todo, err error) {
	err = r.repository.Where("id = ?", todoID).Find(&todo).Error
	return
}

func (r TodoRepository) ReadAll(username string) (todos []todoDomains.Todo, err error) {
	err = r.repository.Where("username = ?", username).Find(&todos).Error
	return
}

func (r TodoRepository) UpdateTodo(updatedTodo todoDomains.Todo) (err error) {
	err = r.repository.Save(&updatedTodo).Error
	return
}

func (r TodoRepository) DeleteTodo(todoID uint) (err error) {
	var todo todoDomains.Todo
	err = r.repository.Where("id = ?", todoID).Delete(&todo).Error
	return
}
