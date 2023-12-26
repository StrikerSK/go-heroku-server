package userHandlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/src/responses"
	"go-heroku-server/api/todo/domain"
	todoPorts "go-heroku-server/api/todo/ports"
	userHandlers "go-heroku-server/api/user/handler"
	"log"
	"net/http"
	"strconv"
)

const (
	todoIdContextKey = "todo_ID"
)

type TodoHandler struct {
	userMiddleware  userHandlers.UserAuthMiddleware
	todoService     todoPorts.ITodoService
	responseService responses.ResponseFactory
}

func NewTodoHandler(userMiddleware userHandlers.UserAuthMiddleware, todoService todoPorts.ITodoService, responseService responses.ResponseFactory) TodoHandler {
	return TodoHandler{
		userMiddleware:  userMiddleware,
		todoService:     todoService,
		responseService: responseService,
	}
}

func (h TodoHandler) EnrichRouter(router *mux.Router) {
	todoRoute := router.PathPrefix("/todo").Subrouter()
	todoRoute.Use(h.userMiddleware.VerifyToken)
	todoRoute.Handle("", http.HandlerFunc(h.createTodo)).Methods(http.MethodPost)
	todoRoute.Handle("/{id}", http.HandlerFunc(h.readTodo)).Methods(http.MethodGet)
	todoRoute.Handle("/{id}", http.HandlerFunc(h.updateTodo)).Methods(http.MethodPut)
	todoRoute.Handle("/{id}", http.HandlerFunc(h.deleteTodo)).Methods(http.MethodDelete)

	todosRoute := router.PathPrefix("/todos").Subrouter()
	todosRoute.Use(h.userMiddleware.VerifyToken)
	todosRoute.Handle("", http.HandlerFunc(h.readTodos)).Methods(http.MethodGet)
}

func ResolveTodoID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), todoIdContextKey, uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todo, err := h.parseTodoBody(r)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todo.Username = username
	err = h.todoService.CreateTodo(todo)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	h.responseService.CreateResponse(nil).WriteResponse(w)
	return
}

func (h TodoHandler) readTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todoID, err := h.resolveTodoID(r)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todo, err := h.todoService.ReadTodo(todoID, userID)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	} else {
		h.responseService.CreateResponse(todo).WriteResponse(w)
		return
	}
}

func (h TodoHandler) readTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todos, err := h.todoService.ReadTodos(userID)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	log.Printf("Todos read: success\n")
	h.responseService.CreateResponse(todos).WriteResponse(w)
	return
}

func (h TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := h.resolveTodoID(r)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	userID, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todo, err := h.parseTodoBody(r)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	todo.Id = todoID

	err = h.todoService.UpdateTodo(userID, todo)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	h.responseService.CreateResponse(nil).WriteResponse(w)
	return
}

func (h TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := h.resolveTodoID(r)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	userID, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	err = h.todoService.DeleteTodo(todoID, userID)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	} else {
		h.responseService.CreateResponse(nil).WriteResponse(w)
		return
	}
}

func (TodoHandler) resolveTodoID(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		return 0, errors.NewParseError(err.Error())
	}

	return uint(id), err
}

func (TodoHandler) parseTodoBody(r *http.Request) (todoDomains.Todo, error) {
	var todo todoDomains.Todo
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&todo); err != nil {
		log.Printf("Resolve Todo: %s\n", err.Error())
		return todoDomains.Todo{}, errors.NewParseError(err.Error())
	}

	return todo, nil
}
