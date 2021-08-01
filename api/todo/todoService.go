package todo

import (
	"errors"
	"go-heroku-server/api/src"
	"go-heroku-server/api/user"
	"log"
	"net/http"
	"strconv"
)

func findAllTodos(userID uint) ([]Todo, *src.RequestError) {
	var todos []Todo

	todos, err := readAll(userID)
	if err != nil {
		return nil, &src.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	return todos, nil
}

func getTodo(todoID uint, userID uint) (*Todo, *src.RequestError) {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf(err.Error() + " for id: " + strconv.Itoa(int(todoID)))
		return nil, &src.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if persistedTodo.UserID != userID {
		errorOutput := errors.New("todo access is denied")
		log.Printf(errorOutput.Error() + " for id: " + strconv.Itoa(int(todoID)))
		return nil, &src.RequestError{
			StatusCode: http.StatusForbidden,
			Err:        errorOutput,
		}
	}

	return &persistedTodo, nil
}

func addTodo(userID uint, todo Todo) {
	todo.UserID = userID
	createTodo(todo)
}

func removeTodo(userID, todoID uint) *src.RequestError {
	persistedFile, err := readTodo(todoID)
	if err != nil {
		log.Printf(err.Error() + " for id: " + strconv.Itoa(int(todoID)))
		return nil
	}

	if persistedFile.UserID != userID {
		return &src.RequestError{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("logged user cannot remove todo"),
		}
	}

	log.Printf("Deleted todo with ID: %d", persistedFile.Id)
	if _, err = deleteTodo(persistedFile.Id); err != nil {
		log.Print(err)
		return nil
	}

	return nil
}

func editTodo(userID, todoID uint, updatedTodo Todo) *src.RequestError {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Print(err)
		return &src.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if persistedTodo.UserID != userID {
		outputError := errors.New("user were accessing unowned todo")
		log.Print(outputError)
		return &src.RequestError{
			StatusCode: http.StatusForbidden,
			Err:        outputError,
		}
	}

	updatedTodo.Id = todoID
	if err = updateTodo(updatedTodo); err != nil {
		log.Print(err)
		return &src.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	return nil
}

func markDone(w http.ResponseWriter, r *http.Request) {
	claimId := uint(r.Context().Value(todoIdContextKey).(int64))
	userID, _ := user.ResolveUserContext(r.Context())

	if persistedTodo, err := readTodo(claimId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	} else {
		if persistedTodo.UserID != userID {
			w.WriteHeader(http.StatusForbidden)
			log.Print("user is accessing forbidden file")
			return
		}

		persistedTodo.Done = !persistedTodo.Done
		err = updateTodo(persistedTodo)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}
	}
}
