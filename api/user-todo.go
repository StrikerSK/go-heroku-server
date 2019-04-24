package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"go-heroku-server/config"
)

type UserNumber struct {
	Id     uint   `json:"id"`
	Number uint   `json:"usernumber"`
	Todos  []Todo `json:"todos" gorm:"foreignkey:StudentId"`
}

type Todo struct {
	Id          uint   `json:"id"`
	StudentId   uint   `json:"-"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func GetAllStudents(w http.ResponseWriter, r *http.Request) {

	var users []UserNumber

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Find(&users)

	for index, user := range users {

		db.Model(&user).Related(&user.Todos, "student_id")
		users[index] = user

	}

	json.NewEncoder(w).Encode(users)

	log.Println("Retrieved User!")
}

func GetStudent(w http.ResponseWriter, r *http.Request) {

	var userNumber UserNumber
	var todo []Todo

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.First(&userNumber).Find(&todo)

	userNumber.Todos = todo
	json.NewEncoder(w).Encode(userNumber)

	log.Println("Retrieved User!")
}

func GetStudentTodos(w http.ResponseWriter, r *http.Request) {

	var todos []Todo
	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	vars := mux.Vars(r)
	uri := vars["id"]

	db.Where("student_id = ?", uri).Find(&todos)
	json.NewEncoder(w).Encode(todos)

}

func AddStudentTodo(w http.ResponseWriter, r *http.Request) {

	var student UserNumber
	var todo Todo

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	vars := mux.Vars(r)
	uri := vars["id"]

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&todo)

	newUint, err := strconv.ParseUint(uri, 10, 32)

	if err != nil {
		log.Fatal(err)
	}

	userId := uint(newUint)

	db.Where("id = ?", uri).Find(&student)
	todo.StudentId = userId

	if student.Id == userId {
		db.NewRecord(todo)
		db.Create(&todo)
	}

}
