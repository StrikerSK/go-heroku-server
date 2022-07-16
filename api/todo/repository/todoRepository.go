package todoRepositories

import (
	"github.com/jinzhu/gorm"
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

func (r TodoRepository) CreateTodo(todo todoDomains.Todo) error {
	return r.db.Create(&todo).Error
}

func (r TodoRepository) ReadTodo(todoID uint) (todo todoDomains.Todo, err error) {
	err = r.db.Where("id = ?", todoID).Find(&todo).Error
	return
}

func (r TodoRepository) ReadTodos(username string) (todos []todoDomains.Todo, err error) {
	err = r.db.Where("username = ?", username).Find(&todos).Error
	return
}

func (r TodoRepository) UpdateTodo(todo todoDomains.Todo) (err error) {
	err = r.db.Save(&todo).Error
	return
}

func (r TodoRepository) DeleteTodo(todo todoDomains.Todo) (err error) {
	err = r.db.Delete(&todo).Error
	return
}
