package todo

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/responses"
	"go-heroku-server/config"
	"log"
	"net/http"
	"strconv"
)

const (
	todoBodyContextKey = "todo_body"
	todoIdContextKey   = "todo_ID"
)

func EnrichRouteWithTodo(router *mux.Router) {

	config.InitializeType("Todo", &Todo{})

	subroute := router.PathPrefix("/todo").Subrouter()
	subroute.Handle("/add", handler.VerifyJwtToken(ResolveTodo(http.HandlerFunc(controllerAddTodo)))).Methods(http.MethodPost)
	subroute.Handle("/{id}", handler.VerifyJwtToken(ResolveTodoID(http.HandlerFunc(controllerGetTodo)))).Methods(http.MethodGet)
	subroute.Handle("/{id}", handler.VerifyJwtToken(ResolveTodoID(ResolveTodo(http.HandlerFunc(controllerEditTodo))))).Methods(http.MethodPut)
	subroute.Handle("/{id}", ResolveTodoID(handler.VerifyJwtToken(http.HandlerFunc(controllerRemoveTodo)))).Methods(http.MethodDelete)
	subroute.Handle("/{id}/done", handler.VerifyJwtToken(ResolveTodoID(http.HandlerFunc(markDone)))).Methods(http.MethodPost, http.MethodGet, http.MethodPut)

	todosSubroute := router.PathPrefix("/todos").Subrouter()
	todosSubroute.Handle("/", handler.VerifyJwtToken(http.HandlerFunc(controllerFindAllTodos))).Methods(http.MethodGet)

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
		var todo Todo
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

func controllerFindAllTodos(w http.ResponseWriter, r *http.Request) {
	userID, _ := handler.ResolveUserContext(r.Context())
	findAllTodos(userID).WriteResponse(w)
}

func controllerGetTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := handler.ResolveUserContext(r.Context())
	todoID := resolveTodoID(r.Context())
	getTodo(todoID, userID).WriteResponse(w)
}

func controllerAddTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := handler.ResolveUserContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(Todo)
	todo.UserID = userID
	addTodo(todo).WriteResponse(w)
}

func controllerRemoveTodo(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := handler.ResolveUserContext(r.Context())
	removeTodo(userID, todoID).WriteResponse(w)
}

func controllerEditTodo(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := handler.ResolveUserContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(Todo)
	editTodo(userID, todoID, todo).WriteResponse(w)
}

func resolveTodoID(context context.Context) uint {
	return uint(context.Value(todoIdContextKey).(int64))
}
