package todo

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/user"
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
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), todoBodyContextKey, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func controllerFindAllTodos(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	if todos, requestError := findAllTodos(userID); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(todos)
	}
}

func controllerGetTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	todoID := uint(r.Context().Value(todoIdContextKey).(int64))

	if persistedTodo, requestError := getTodo(todoID, userID); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(persistedTodo)
		_, _ = w.Write(payload)
	}
}

func controllerAddTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(Todo)
	addTodo(userID, todo)
}

func controllerRemoveTodo(w http.ResponseWriter, r *http.Request) {
	todoID := uint(r.Context().Value(todoIdContextKey).(int64))
	userID, _ := user.ResolveUserContext(r.Context())

	if requestError := removeTodo(userID, todoID); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func controllerEditTodo(w http.ResponseWriter, r *http.Request) {
	todoID := uint(r.Context().Value(todoIdContextKey).(int64))
	userID, _ := user.ResolveUserContext(r.Context())
	todo := r.Context().Value(todoBodyContextKey).(Todo)

	if requestError := editTodo(userID, todoID, todo); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
