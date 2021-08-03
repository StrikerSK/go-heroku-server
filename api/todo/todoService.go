package todo

import (
	"errors"
	"go-heroku-server/api/src"
	"go-heroku-server/api/user"
	"log"
	"net/http"
)

func findAllTodos(userID uint) src.IResponse {
	var todos []Todo

	todos, err := readAll(userID)
	if err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	}

	return src.NewResponse(todos)
}

func getTodo(todoID uint, userID uint) src.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		return src.NewErrorResponse(http.StatusNotFound, err)
	}

	if persistedTodo.UserID != userID {
		return src.NewErrorResponse(http.StatusForbidden, errors.New("todo access is denied"))
	}

	return src.NewResponse(persistedTodo)
}

func addTodo(todo Todo) src.IResponse {
	createTodo(todo)
	return src.NewEmptyResponse(http.StatusCreated)
}

func removeTodo(userID, todoID uint) src.IResponse {
	persistedFile, err := readTodo(todoID)
	if err != nil {
		log.Printf("%s for todo id: %d\n", err.Error(), int(todoID))
		return src.NewEmptyResponse(http.StatusOK)
	}

	if persistedFile.UserID != userID {
		return src.NewErrorResponse(http.StatusForbidden, errors.New("access denied"))
	}

	if err = deleteTodo(persistedFile.Id); err != nil {
		log.Print(err)
		return src.NewEmptyResponse(http.StatusOK)
	}

	log.Printf("Deleted Todo with ID: %d", persistedFile.Id)
	return src.NewEmptyResponse(http.StatusOK)
}

func editTodo(userID, todoID uint, updatedTodo Todo) src.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] not found\n", todoID)
		return src.NewEmptyResponse(http.StatusNotFound)
	}

	if persistedTodo.UserID != userID {
		outputError := errors.New("access denied")
		return src.NewErrorResponse(http.StatusForbidden, outputError)
	}

	updatedTodo.Id = todoID
	updatedTodo.UserID = persistedTodo.UserID

	if err = updateTodo(updatedTodo); err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	}

	return src.EmptyResponse(http.StatusOK)
}

func markDone(w http.ResponseWriter, r *http.Request) {
	claimId := resolveTodoID(r.Context())
	userID, _ := user.ResolveUserContext(r.Context())

	persistedTodo, err := readTodo(claimId)
	if err != nil {
		src.NewErrorResponse(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

	if persistedTodo.UserID != userID {
		err = errors.New("access denied")
		src.NewErrorResponse(http.StatusForbidden, err).WriteResponse(w)
		return
	}

	persistedTodo.Done = !persistedTodo.Done
	err = updateTodo(persistedTodo)
	if err != nil {
		src.NewErrorResponse(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

	src.NewResponse(http.StatusOK).WriteResponse(w)
}
