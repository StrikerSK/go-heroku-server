package userHandlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/responses"
	"go-heroku-server/api/todo/domain"
	todoPorts "go-heroku-server/api/todo/ports"
	userHandlers "go-heroku-server/api/user/handler"
	"go-heroku-server/config"
	"log"
	"net/http"
	"strconv"
)

const (
	todoBodyContextKey = "todo_body"
	todoIdContextKey   = "todo_ID"
)

type TodoHandler struct {
	userMiddleware userHandlers.UserAuthMiddleware
	todoService    todoPorts.ITodoService
}

func NewTodoHandler(userMiddleware userHandlers.UserAuthMiddleware, todoService todoPorts.ITodoService) TodoHandler {
	return TodoHandler{
		userMiddleware: userMiddleware,
		todoService:    todoService,
	}
}

func (h TodoHandler) EnrichRouter(router *mux.Router) {
	config.InitializeType("Todo", &todoDomains.Todo{})

	todoRoute := router.PathPrefix("/todo").Subrouter()
	todoRoute.Handle("", h.userMiddleware.VerifyToken(ResolveTodo(http.HandlerFunc(h.createTodo)))).Methods(http.MethodPost)
	todoRoute.Handle("/{id}", h.userMiddleware.VerifyToken(ResolveTodoID(http.HandlerFunc(h.readTodo)))).Methods(http.MethodGet)
	todoRoute.Handle("/{id}", h.userMiddleware.VerifyToken(ResolveTodoID(ResolveTodo(http.HandlerFunc(h.updateTodo))))).Methods(http.MethodPut)
	todoRoute.Handle("/{id}", ResolveTodoID(h.userMiddleware.VerifyToken(http.HandlerFunc(h.deleteTodo)))).Methods(http.MethodDelete)
	//todoRoute.Handle("/{id}/done", h.userMiddleware.VerifyToken(ResolveTodoID(http.HandlerFunc(markDone)))).Methods(http.MethodPost, http.MethodGet, http.MethodPut)

	todosRoute := router.PathPrefix("/todos").Subrouter()
	todosRoute.Handle("", h.userMiddleware.VerifyToken(http.HandlerFunc(h.readTodos))).Methods(http.MethodGet)
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

func ResolveTodo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var todo todoDomains.Todo
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&todo); err != nil {
			log.Printf("Resolve Todo: %s\n", err.Error())
			responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), todoBodyContextKey, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	username, _ := h.userMiddleware.GetUsernameFromContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(todoDomains.Todo)

	todo.Username = username
	if err := h.todoService.CreateTodo(todo); err != nil {
		log.Printf("Todo create: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusCreated, nil).WriteResponse(w)
		return
	}
}

func (h TodoHandler) readTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := h.userMiddleware.GetUsernameFromContext(r.Context())
	todoID := resolveTodoID(r.Context())
	todo, err := h.todoService.ReadTodo(todoID, userID)
	if err != nil {
		log.Printf("Todo create: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, todo).WriteResponse(w)
		return
	}
}

func (h TodoHandler) readTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		log.Printf("Todos read: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	todos, err := h.todoService.ReadTodos(userID)
	if err != nil {
		log.Printf("Todos read: %s\n", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	log.Printf("Todos read: success\n")
	responses.CreateResponse(http.StatusOK, todos).WriteResponse(w)
	return
}

func (h TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := h.userMiddleware.GetUsernameFromContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(todoDomains.Todo)

	if err := h.todoService.UpdateTodo(todoID, userID, todo); err != nil {
		log.Printf("Todo [%d] edit: %v\n", todoID, err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	log.Printf("Todo [%d] edit: success\n", todoID)
	responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
	return
}

func (h TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := h.userMiddleware.GetUsernameFromContext(r.Context())

	err := h.todoService.DeleteTodo(todoID, userID)
	if err != nil {
		log.Printf("Todo [%d] delete: %s\n", todoID, err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
		return
	}
}

func resolveTodoID(context context.Context) uint {
	return uint(context.Value(todoIdContextKey).(int64))
}
