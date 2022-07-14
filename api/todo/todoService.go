package todo

import (
	"go-heroku-server/api/src/responses"
	"log"
	"net/http"
)

func findAllTodos(userID uint) responses.IResponse {
	todos, err := readAll(userID)
	if err != nil {
		log.Printf("Todos read: %s\n", err.Error())
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	log.Printf("Todos read: success\n")
	return responses.CreateResponse(http.StatusOK, todos)
}

func getTodo(todoID uint, userID uint) responses.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] read: %s\n", todoID, err.Error())
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	if persistedTodo.UserID != userID {
		log.Printf("Todo [%d] read: access denied\n", todoID)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	log.Printf("Todo [%d] read: success\n", todoID)
	return responses.CreateResponse(http.StatusOK, persistedTodo)
}

func addTodo(todo Todo) responses.IResponse {
	createTodo(todo)
	return responses.CreateResponse(http.StatusCreated, nil)
}

func removeTodo(userID, todoID uint) responses.IResponse {
	persistedFile, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] delete: %s\n", todoID, err.Error())
		return responses.CreateResponse(http.StatusOK, nil)
	}

	if persistedFile.UserID != userID {
		log.Printf("Todo [%d] delete: access denied\n", todoID)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	if err = deleteTodo(persistedFile.Id); err != nil {
		log.Printf("Todo [%d] delete: %s\n", todoID, err.Error())
		return responses.CreateResponse(http.StatusOK, nil)
	}

	log.Printf("Todo [%d] delete: success\n", todoID)
	return responses.CreateResponse(http.StatusOK, nil)
}

func editTodo(userID, todoID uint, updatedTodo Todo) responses.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] edit: %s\n", todoID, err.Error())
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	if persistedTodo.UserID != userID {
		log.Printf("Todo [%d] edit: access denied\n", todoID)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	updatedTodo.Id = todoID
	updatedTodo.UserID = persistedTodo.UserID

	if err = updateTodo(updatedTodo); err != nil {
		log.Printf("Todo [%d] edit: %s\n", todoID, err.Error())
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	log.Printf("Todo [%d] edit: success\n", todoID)
	return responses.CreateResponse(http.StatusOK, nil)
}

func markDone(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := handler.ResolveUserContext(r.Context())

	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] toggle: %s\n", todoID, err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	persistedTodo.Done = !persistedTodo.Done
	editTodo(userID, todoID, persistedTodo)

	responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
}
