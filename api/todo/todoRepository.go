package todo

import (
	"go-heroku-server/config"
)

func createTodo(newTodo Todo) {
	config.GetDatabaseInstance().NewRecord(newTodo)
	config.GetDatabaseInstance().Create(&newTodo)
}

func readTodo(todoID uint) (todo Todo, err error) {
	err = config.GetDatabaseInstance().Where("id = ?", todoID).Find(&todo).Error
	return
}

func readAll(userID uint) (todos []Todo, err error) {
	err = config.GetDatabaseInstance().Where("user_id = ?", userID).Find(&todos).Error
	return
}

func updateTodo(updatedTodo Todo) (err error) {
	err = config.GetDatabaseInstance().Save(&updatedTodo).Error
	return
}

func deleteTodo(fileID interface{}) (err error) {
	var t Todo
	err = config.GetDatabaseInstance().Where("id = ?", fileID).Delete(&t).Error
	return
}
