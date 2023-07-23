package main

import (
	"fmt"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	fileHandlers "go-heroku-server/api/files/v1/handler"
	fileRepositories "go-heroku-server/api/files/v1/repository"
	fileServices "go-heroku-server/api/files/v1/service"
	locationHandlers "go-heroku-server/api/location/handler"
	locationRepositories "go-heroku-server/api/location/repository"
	locationServices "go-heroku-server/api/location/service"
	"go-heroku-server/api/src/responses"
	todoHandlers "go-heroku-server/api/todo/handler"
	todoRepositories "go-heroku-server/api/todo/repository"
	todoServices "go-heroku-server/api/todo/service"
	userHandlers "go-heroku-server/api/user/handler"
	userRepositories "go-heroku-server/api/user/repository"
	userServices "go-heroku-server/api/user/service"
	"html/template"
	"net/http"
	"os"

	"go-heroku-server/config"
)

func serveMainPage(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("static/index.html"))
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func init() {
	config.GetCacheInstance()
}

// Go application entrypoint
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "8080"
	}

	responseService := responses.NewResponseFactory()
	databaseInstance := config.GetDatabaseInstance()

	userRepository := userRepositories.NewUserRepository(databaseInstance)
	userService := userServices.NewUserService(userRepository)

	userTokenService := userServices.NewTokenService("Wow, much safe", 3600)
	userMiddleware := userHandlers.NewUserAuthMiddleware(userTokenService, responseService)
	userHdl := userHandlers.NewUserHandler(userService, userMiddleware, userTokenService, responseService)

	todoRepo := todoRepositories.NewTodoRepository(databaseInstance)
	todoService := todoServices.NewTodoService(todoRepo)
	todoHdl := todoHandlers.NewTodoHandler(userMiddleware, todoService, responseService)

	fileRepo := fileRepositories.NewFileDatabaseRepository(databaseInstance)
	fileSrv := fileServices.NewFileService(fileRepo)
	fileHdl := fileHandlers.NewMuxFileHandler(fileSrv, userMiddleware, responseService)

	locationRepo := locationRepositories.NewLocationRepository(databaseInstance)
	locationSrv := locationServices.NewLocationService(locationRepo)
	locationHdl := locationHandlers.NewLocationHandler(locationSrv, userMiddleware, responseService)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", serveMainPage)

	userHdl.EnrichRouter(router)
	todoHdl.EnrichRouter(router)
	fileHdl.EnrichRouter(router)
	locationHdl.EnrichRouter(router)

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Listening on port ", port)
	fmt.Println(http.ListenAndServe(":"+port, handler))
}
