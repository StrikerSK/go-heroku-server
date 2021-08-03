package todo

import "go-heroku-server/config"

func createTodo(newTodo Todo) {
	config.DBConnection.NewRecord(newTodo)
	config.DBConnection.Create(&newTodo)
}

func readTodo(todoID uint) (todo Todo, err error) {
	err = config.DBConnection.Where("id = ?", todoID).Find(&todo).Error
	return
}

func readAll(userID uint) (todos []Todo, err error) {
	err = config.DBConnection.Where("user_id = ?", userID).Find(&todos).Error
	return
}

func updateTodo(updatedTodo Todo) (err error) {
	err = config.DBConnection.Save(&updatedTodo).Error
	return
}

func deleteTodo(fileID interface{}) (err error) {
	var t Todo
	err = config.DBConnection.Where("id = ?", fileID).Delete(&t).Error
	return
}
