package todo

import (
	"encoding/json"
	"log"
	"net/http"

	"go-heroku-server/config"
)

type UserNumber struct {
	Id     uint   `json:"id"`
	Number uint   `json:"userNumber"`
	Todos  []Todo `json:"todos" gorm:"foreignkey:StudentId"`
}

func GetAllStudents(w http.ResponseWriter, r *http.Request) {

	var users []UserNumber
	config.DBConnection.Find(&users)
	for index, user := range users {
		config.DBConnection.Model(&user).Related(&user.Todos, "student_id")
		users[index] = user
	}

	json.NewEncoder(w).Encode(users)

	log.Println("Retrieved User!")
}

func GetStudent(w http.ResponseWriter, r *http.Request) {

	var userNumber UserNumber
	var todo []Todo

	config.DBConnection.First(&userNumber).Find(&todo)

	userNumber.Todos = todo
	json.NewEncoder(w).Encode(userNumber)

	log.Println("Retrieved User!")
}
