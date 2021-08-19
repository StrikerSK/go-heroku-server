package todo

import (
	"go-heroku-server/api/src/responses"
	"go-heroku-server/api/user"
	"log"
	"net/http"
)

func findAllTodos(userID uint) responses.IResponse {
	todos, err := readAll(userID)
	if err != nil {
		log.Printf("Todos read: %s\n", err.Error())
		return responses.NewEmptyResponse(http.StatusBadRequest)
	}

	log.Printf("Todos read: success\n")
	return responses.NewResponse(todos)
}

func getTodo(todoID uint, userID uint) responses.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] read: %s\n", todoID, err.Error())
		return responses.NewEmptyResponse(http.StatusNotFound)
	}

	if persistedTodo.UserID != userID {
		log.Printf("Todo [%d] read: access denied\n", todoID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	log.Printf("Todo [%d] read: success\n", todoID)
	return responses.NewResponse(persistedTodo)
}

func addTodo(todo Todo) responses.IResponse {
	createTodo(todo)
	return responses.NewEmptyResponse(http.StatusCreated)
}

func removeTodo(userID, todoID uint) responses.IResponse {
	persistedFile, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] delete: %s\n", todoID, err.Error())
		return responses.NewEmptyResponse(http.StatusOK)
	}

	if persistedFile.UserID != userID {
		log.Printf("Todo [%d] delete: access denied\n", todoID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	if err = deleteTodo(persistedFile.Id); err != nil {
		log.Printf("Todo [%d] delete: %s\n", todoID, err.Error())
		return responses.NewEmptyResponse(http.StatusOK)
	}

	log.Printf("Todo [%d] delete: success\n", todoID)
	return responses.NewEmptyResponse(http.StatusOK)
}

func editTodo(userID, todoID uint, updatedTodo Todo) responses.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] edit: %s\n", todoID, err.Error())
		return responses.NewEmptyResponse(http.StatusNotFound)
	}

	if persistedTodo.UserID != userID {
		log.Printf("Todo [%d] edit: access denied\n", todoID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	updatedTodo.Id = todoID
	updatedTodo.UserID = persistedTodo.UserID

	if err = updateTodo(updatedTodo); err != nil {
		log.Printf("Todo [%d] edit: %s\n", todoID, err.Error())
		return responses.NewEmptyResponse(http.StatusBadRequest)
	}

	log.Printf("Todo [%d] edit: success\n", todoID)
	return responses.NewEmptyResponse(http.StatusOK)
}

func markDone(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := user.ResolveUserContext(r.Context())

	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] toggle: %s\n", todoID, err.Error())
		responses.NewErrorResponse(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

	persistedTodo.Done = !persistedTodo.Done
	editTodo(userID, todoID, persistedTodo)

	responses.NewResponse(http.StatusOK).WriteResponse(w)
}
