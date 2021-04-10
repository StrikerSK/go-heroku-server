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

func updateTodo(updatedTodo Todo) (err error) {
	err = config.DBConnection.Model(&Todo{}).Update(&updatedTodo).Error
	return
}

func deleteTodo(fileID interface{}) (todo Todo, err error) {
	err = config.DBConnection.Where("id = ?", fileID).Delete(&todo).Error
	return
}
