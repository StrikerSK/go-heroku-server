package todoPorts

import (
	todoDomains "go-heroku-server/api/todo/domain"
)

type ITodoRepository interface {
	CreateTodo(todoDomains.Todo) error
	ReadTodo(uint) (todoDomains.Todo, error)
	ReadTodos(string) ([]todoDomains.Todo, error)
	UpdateTodo(todoDomains.Todo) error
	DeleteTodo(todo todoDomains.Todo) (err error)
}
