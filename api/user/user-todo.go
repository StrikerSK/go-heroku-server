package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"go-heroku-server/config"
)

type Todo struct {
	Id          uint   `json:"id"`
	StudentId   uint   `json:"-"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func GetStudentTodos(w http.ResponseWriter, r *http.Request) {

	var todos []Todo
	uri := mux.Vars(r)["id"]
	config.DBConnection.Where("student_id = ?", uri).Find(&todos)
	_ = json.NewEncoder(w).Encode(todos)

}

func AddStudentTodo(w http.ResponseWriter, r *http.Request) {

	var student UserNumber
	var todo Todo

	vars := mux.Vars(r)
	uri := vars["id"]

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&todo)

	newUint, err := strconv.ParseUint(uri, 10, 32)

	if err != nil {
		log.Fatal(err)
	}

	userId := uint(newUint)

	config.DBConnection.Where("id = ?", uri).Find(&student)
	todo.StudentId = userId

	if student.Id == userId {
		config.DBConnection.NewRecord(todo)
		config.DBConnection.Create(&todo)
	}

}
