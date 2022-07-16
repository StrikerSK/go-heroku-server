package todoPorts

import (
	todoDomains "go-heroku-server/api/todo/domain"
)

type ITodoRepository interface {
	CreateTodo(todoDomains.Todo) error
	ReadTodo(uint) (todoDomains.Todo, error)
	ReadAll(string) ([]todoDomains.Todo, error)
	UpdateTodo(todoDomains.Todo) error
	DeleteTodo(uint) (err error)
}
