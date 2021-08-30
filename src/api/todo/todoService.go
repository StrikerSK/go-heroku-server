package todo

import (
	"go-heroku-server/src/api/user"
	"go-heroku-server/src/responses"
	"log"
	"net/http"
)

func findAllTodos(userID uint) responses.IResponse {
	todos, err := readAll(userID)
	if err != nil {
		log.Printf("Todos read: %v\n", err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	log.Printf("Todos read: success\n")
	return responses.CreateResponse(http.StatusOK, todos)
}

func getTodo(todoID uint, userID uint) responses.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] read: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	if err = persistedTodo.validateAccess(userID); err != nil {
		log.Printf("Todo [%d] read: %v\n", todoID, err)
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
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] delete: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusOK, nil)
	}

	if err = persistedTodo.validateAccess(userID); err != nil {
		log.Printf("Todo [%d] delete: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	if err = deleteTodo(persistedTodo.Id); err != nil {
		log.Printf("Todo [%d] delete: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusOK, nil)
	}

	log.Printf("Todo [%d] delete: success\n", todoID)
	return responses.CreateResponse(http.StatusOK, nil)
}

func editTodo(userID, todoID uint, updatedTodo Todo) responses.IResponse {
	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] edit: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	if err = persistedTodo.validateAccess(userID); err != nil {
		log.Printf("Todo [%d] edit: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	updatedTodo.Id = todoID
	updatedTodo.UserID = persistedTodo.UserID

	if err = updateTodo(updatedTodo); err != nil {
		log.Printf("Todo [%d] edit: %v\n", todoID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	log.Printf("Todo [%d] edit: success\n", todoID)
	return responses.CreateResponse(http.StatusOK, nil)
}

func markDone(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := user.ResolveUserContext(r.Context())

	persistedTodo, err := readTodo(todoID)
	if err != nil {
		log.Printf("Todo [%d] toggle: %v\n", todoID, err)
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	persistedTodo.Done = !persistedTodo.Done
	editTodo(userID, todoID, persistedTodo)

	responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
}
