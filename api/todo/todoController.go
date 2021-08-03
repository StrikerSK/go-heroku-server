package todo

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"net/http"
	"strconv"
)

const (
	todoBodyContextKey = "todo_body"
	todoIdContextKey   = "todo_ID"
)

func EnrichRouteWithTodo(router *mux.Router) {

	config.DBConnection.AutoMigrate(&Todo{})

	subroute := router.PathPrefix("/todo").Subrouter()
	subroute.Handle("/add", user.VerifyJwtToken(ResolveTodo(http.HandlerFunc(controllerAddTodo)))).Methods("POST")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveTodoID(http.HandlerFunc(controllerGetTodo)))).Methods("GET")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveTodoID(ResolveTodo(http.HandlerFunc(controllerEditTodo))))).Methods("PUT")
	subroute.Handle("/{id}", ResolveTodoID(user.VerifyJwtToken(http.HandlerFunc(controllerRemoveTodo)))).Methods("DELETE")
	subroute.Handle("/{id}/done", user.VerifyJwtToken(ResolveTodoID(http.HandlerFunc(markDone)))).Methods("POST", "GET", "PUT")

	todosSubroute := router.PathPrefix("/todos").Subrouter()
	todosSubroute.Handle("/", user.VerifyJwtToken(http.HandlerFunc(controllerFindAllTodos))).Methods("GET")

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
		err := decoder.Decode(&todo)
		if err != nil {
			src.NewErrorResponse(http.StatusInternalServerError, err).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), todoBodyContextKey, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func controllerFindAllTodos(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	findAllTodos(userID).WriteResponse(w)
}

func controllerGetTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	todoID := resolveTodoID(r.Context())
	getTodo(todoID, userID).WriteResponse(w)
}

func controllerAddTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(Todo)
	todo.UserID = userID
	addTodo(todo).WriteResponse(w)
}

func controllerRemoveTodo(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := user.ResolveUserContext(r.Context())
	removeTodo(userID, todoID).WriteResponse(w)
}

func controllerEditTodo(w http.ResponseWriter, r *http.Request) {
	todoID := resolveTodoID(r.Context())
	userID, _ := user.ResolveUserContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(Todo)
	editTodo(userID, todoID, todo).WriteResponse(w)
}

func resolveTodoID(context context.Context) uint {
	return uint(context.Value(todoIdContextKey).(int64))
}
