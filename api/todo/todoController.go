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
	subroute.Handle("/add", user.VerifyJwtToken(ResolveTodo(http.HandlerFunc(addTodo)))).Methods("POST")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveTodoID(http.HandlerFunc(getTodo)))).Methods("GET")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveTodoID(ResolveTodo(http.HandlerFunc(editTodo))))).Methods("PUT")
	subroute.Handle("/{id}", ResolveTodoID(user.VerifyJwtToken(http.HandlerFunc(removeTodo)))).Methods("DELETE")
	subroute.Handle("/{id}/done", user.VerifyJwtToken(ResolveTodoID(http.HandlerFunc(markDone)))).Methods("POST", "GET", "PUT")

	todosSubroute := router.PathPrefix("/todos").Subrouter()
	todosSubroute.Handle("/", user.VerifyJwtToken(http.HandlerFunc(findAllTodos))).Methods("GET")

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
