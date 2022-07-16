package todoPorts

import (
	todoDomains "go-heroku-server/api/todo/domain"
)

type ITodoService interface {
	CreateTodo(todoDomains.Todo) error
	ReadTodo(uint, string) (todoDomains.Todo, error)
	ReadTodos(string) ([]todoDomains.Todo, error)
	UpdateTodo(uint, string, todoDomains.Todo) error
	DeleteTodo(uint, string) error
}
