package todoRepositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/todo/domain"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	db.AutoMigrate(&todoDomains.Todo{})
	return TodoRepository{
		db: db,
	}
}

func (r TodoRepository) CreateTodo(todo todoDomains.Todo) (err error) {
	if err = r.db.Create(&todo).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r TodoRepository) ReadTodo(todoID uint) (todo todoDomains.Todo, err error) {
	if err = r.db.Where("id = ?", todoID).Find(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return todoDomains.Todo{}, errors.NewNotFoundError(fmt.Sprintf("location [%d] not found", todoID))
		} else {
			return todoDomains.Todo{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r TodoRepository) ReadTodos(username string) (todos []todoDomains.Todo, err error) {
	if err = r.db.Where("username = ?", username).Find(&todos).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r TodoRepository) UpdateTodo(todo todoDomains.Todo) (err error) {
	if err = r.db.Save(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("todo [%d] not found", todo.Id))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r TodoRepository) DeleteTodo(todo todoDomains.Todo) (err error) {
	if err = r.db.Delete(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("todo [%d] not found", todo.Id))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}
