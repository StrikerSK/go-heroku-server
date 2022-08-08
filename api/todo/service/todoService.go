package todoServices

import (
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/todo/domain"
	todoPorts "go-heroku-server/api/todo/ports"
	"log"
)

type TodoService struct {
	repository todoPorts.ITodoRepository
}

func NewTodoService(repository todoPorts.ITodoRepository) TodoService {
	return TodoService{
		repository: repository,
	}
}

func (s TodoService) CreateTodo(todo todoDomains.Todo) error {
	return s.repository.CreateTodo(todo)
}

func (s TodoService) ReadTodo(todoID uint, username string) (todoDomains.Todo, error) {
	todo, err := s.repository.ReadTodo(todoID)
	if err != nil {
		return todoDomains.Todo{}, err
	}

	if username != todo.Username {
		return todoDomains.Todo{}, errors.NewForbiddenError("forbidden access to resource")
	}

	return todo, err
}

func (s TodoService) ReadTodos(username string) ([]todoDomains.Todo, error) {
	return s.repository.ReadTodos(username)
}

func (s TodoService) UpdateTodo(username string, updatedTodo todoDomains.Todo) error {
	if _, err := s.ReadTodo(updatedTodo.Id, username); err != nil {
		return err
	}

	updatedTodo.Username = username

	if err := s.repository.UpdateTodo(updatedTodo); err != nil {
		log.Printf("Todo [%s] edit: %v\n", username, err)
		return err
	}

	return nil
}

func (s TodoService) DeleteTodo(todoID uint, username string) error {
	persistedTodo, err := s.ReadTodo(todoID, username)
	if err != nil {
		return err
	}

	if err = s.repository.DeleteTodo(persistedTodo); err != nil {
		return err
	}

	return nil
}
