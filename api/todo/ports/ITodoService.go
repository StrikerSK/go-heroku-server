package todoPorts

import (
	todoDomains "go-heroku-server/api/todo/domain"
)

type ITodoService interface {
	ReadTodos(string) ([]todoDomains.Todo, error)
	ReadTodo(uint, string) (todoDomains.Todo, error)
	CreateTodo(todoDomains.Todo) error
	DeleteTodo(uint, string) error
	EditTodo(uint, string, todoDomains.Todo) error
}
