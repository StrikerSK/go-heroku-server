package main

import (
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	fileHandlers "go-heroku-server/api/files/v2/handler"
	fileRepositories "go-heroku-server/api/files/v2/repository"
	fileServices "go-heroku-server/api/files/v2/service"
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
	"go-heroku-server/config"
	"go-heroku-server/config/database"
	"html/template"
	"log"
	"net/http"
)

func serveMainPage(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("static/index.html"))
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// Go application entrypoint
func main() {
	viperConfiguration := config.ReadConfiguration()
	applicationConfiguration := viperConfiguration.Application
	databaseConfiguration := viperConfiguration.Database
	authorizationConfiguration := viperConfiguration.Authorization

	responseService := responses.NewResponseFactory()
	databaseInstance := database.CreateDB(databaseConfiguration)

	userRepository := userRepositories.NewUserRepository(databaseInstance)
	userService := userServices.NewUserService(userRepository)

	userTokenService := userServices.NewTokenService(authorizationConfiguration)
	userMiddleware := userHandlers.NewUserAuthMiddleware(userTokenService, responseService, authorizationConfiguration)
	userHdl := userHandlers.NewUserHandler(userService, userMiddleware, userTokenService, responseService)

	todoRepo := todoRepositories.NewTodoRepository(databaseInstance)
	todoService := todoServices.NewTodoService(todoRepo)
	todoHdl := todoHandlers.NewTodoHandler(userMiddleware, todoService, responseService)

	fileRepo := fileRepositories.NewFileDatabaseRepository(databaseInstance)
	fileMetadataRepo := fileRepositories.NewFileMetadataRepository(databaseInstance)
	fileSrv := fileServices.NewFileService(fileMetadataRepo, fileRepo)
	fileHdl := fileHandlers.NewMuxFileHandler(fileSrv, userMiddleware, responseService)

	locationRepo := locationRepositories.NewLocationRepository(databaseInstance)
	locationSrv := locationServices.NewLocationService(locationRepo)
	locationHdl := locationHandlers.NewLocationHandler(locationSrv, userMiddleware, responseService)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", serveMainPage)

	userRoute := router.PathPrefix(applicationConfiguration.ContextPath).Subrouter()
	userHdl.EnrichRouter(userRoute)
	todoHdl.EnrichRouter(userRoute)
	fileHdl.EnrichRouter(userRoute)
	locationHdl.EnrichRouter(userRoute)

	handler := cors.AllowAll().Handler(router)

	log.Println("Listening on port:", applicationConfiguration.Port)
	log.Println("Using context path:", applicationConfiguration.ContextPath)
	log.Println(http.ListenAndServe(":"+applicationConfiguration.Port, handler))
}
