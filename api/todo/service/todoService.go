package todoServices

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-heroku-server/api/todo/domain"
	todoPorts "go-heroku-server/api/todo/ports"
	"go-heroku-server/api/types/errors"
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
		if err == gorm.ErrRecordNotFound {
			return todoDomains.Todo{}, errors.NewNotFoundError(fmt.Sprintf("todo [%d] not found", todoID))
		} else {
			return todoDomains.Todo{}, err
		}
	} else {
		if username != todo.Username {
			return todoDomains.Todo{}, errors.NewForbiddenError("forbidden access to resource")
		}
		return s.repository.ReadTodo(todoID)
	}
}

func (s TodoService) ReadTodos(username string) ([]todoDomains.Todo, error) {
	return s.repository.ReadTodos(username)
}

func (s TodoService) UpdateTodo(todoID uint, username string, updatedTodo todoDomains.Todo) error {
	_, err := s.ReadTodo(todoID, username)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(fmt.Sprintf("todo [%d] not found", todoID))
		} else {
			return err
		}
	} else {
		updatedTodo.Username = username
		updatedTodo.Id = todoID

		if err = s.repository.UpdateTodo(updatedTodo); err != nil {
			log.Printf("Todo [%s] edit: %v\n", username, err)
			return err
		}

		return nil
	}
}

func (s TodoService) DeleteTodo(todoID uint, username string) error {
	persistedFile, err := s.ReadTodo(todoID, username)
	if err != nil {
		return err
	}

	if err = s.repository.DeleteTodo(persistedFile); err != nil {
		return err
	}

	return nil
}
