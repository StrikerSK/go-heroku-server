package todo

import (
	"encoding/json"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"log"
	"net/http"
	"strconv"
)

func findAllTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	claimId := r.Context().Value(user.UserIdContextKey).(uint)
	config.DBConnection.Where("user_id = ?", claimId).Find(&todos)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(user.UserIdContextKey)
	todoID := uint(r.Context().Value(todoIdContextKey).(int64))

	if persistedTodo, err := readTodo(todoID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf(err.Error() + " for id: " + strconv.Itoa(int(todoID)))
		return
	} else {
		if persistedTodo.UserID != userID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(persistedTodo)
		_, _ = w.Write(payload)
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	claimId := r.Context().Value(user.UserIdContextKey).(uint)
	todo := r.Context().Value(todoBodyContextKey).(Todo)
	todo.UserID = claimId
	createTodo(todo)
}

func removeTodo(w http.ResponseWriter, r *http.Request) {
	fileID := uint(r.Context().Value(todoIdContextKey).(int64))
	claimId := r.Context().Value(user.UserIdContextKey)

	if persistedFile, err := readTodo(fileID); err != nil {
		w.WriteHeader(http.StatusOK)
		log.Printf(err.Error() + " for id: " + strconv.Itoa(int(fileID)))
		return
	} else {
		if persistedFile.UserID != claimId {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		log.Printf("Deleted todo with ID: %d", persistedFile.Id)
		_, _ = deleteTodo(persistedFile.Id)
	}
}

func editTodo(w http.ResponseWriter, r *http.Request) {
	todoID := uint(r.Context().Value(todoIdContextKey).(int64))
	userID := r.Context().Value(user.UserIdContextKey).(uint)
	todo := r.Context().Value(todoBodyContextKey).(Todo)

	if persistedTodo, err := readTodo(todoID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Print(err)
		return
	} else {
		if persistedTodo.UserID != userID {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("user were accessing unowned todo")
			return
		}
	}

	todo.Id = todoID
	if err := updateTodo(todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}
}

func markDone(w http.ResponseWriter, r *http.Request) {
	claimId := uint(r.Context().Value(todoIdContextKey).(int64))
	userID := r.Context().Value(user.UserIdContextKey).(uint)

	if persistedTodo, err := readTodo(claimId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	} else {
		if persistedTodo.UserID != userID {
			w.WriteHeader(http.StatusForbidden)
			log.Print("user is accessing forbidden file")
			return
		}

		persistedTodo.Done = !persistedTodo.Done
		err = updateTodo(persistedTodo)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}
	}
}
