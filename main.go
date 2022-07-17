package main

import (
	"fmt"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	fileHandlers "go-heroku-server/api/files/handler"
	fileRepositories "go-heroku-server/api/files/repository"
	fileServices "go-heroku-server/api/files/service"
	"go-heroku-server/api/location/domain"
	locationHandlers "go-heroku-server/api/location/handler"
	"go-heroku-server/api/location/image"
	locationRepositories "go-heroku-server/api/location/repository"
	"go-heroku-server/api/location/restaurant"
	locationServcices "go-heroku-server/api/location/service"
	"go-heroku-server/api/src/responses"
	todoHandlers "go-heroku-server/api/todo/handler"
	todoRepositories "go-heroku-server/api/todo/repository"
	todoServices "go-heroku-server/api/todo/service"
	userAuth "go-heroku-server/api/user/auth"
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
	config.GetDatabaseInstance().AutoMigrate(&locationDomains.UserLocationEntity{}, &image.LocationImage{}, &restaurant.RestaurantLocation{})
	config.GetCacheInstance()
}

//Go application entrypoint
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "4000"
	}

	responseService := responses.NewResponseService()

	userRepository := userRepositories.NewUserRepository(config.GetDatabaseInstance())
	userService := userServices.NewUserService(userRepository)

	userTokenService := userAuth.NewTokenService()
	userMiddleware := userHandlers.NewUserAuthMiddleware(userTokenService)
	userHdl := userHandlers.NewUserHandler(userService, userMiddleware, userTokenService, responseService)

	todoRepo := todoRepositories.NewTodoRepository(config.GetDatabaseInstance())
	todoService := todoServices.NewTodoService(todoRepo)
	todoHdl := todoHandlers.NewTodoHandler(userMiddleware, todoService)

	fileRepo := fileRepositories.NewFileRepository(config.GetDatabaseInstance())
	fileSrv := fileServices.NewFileService(fileRepo)
	fileHdl := fileHandlers.NewFileHandler(fileSrv, userMiddleware)

	locationRepo := locationRepositories.NewLocationRepository(config.GetDatabaseInstance())
	locationSrv := locationServcices.NewLocationService(locationRepo)
	locationHdl := locationHandlers.NewLocationHandler(locationSrv, userMiddleware)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", serveMainPage)

	userHdl.EnrichRouter(router)
	todoHdl.EnrichRouter(router)
	fileHdl.EnrichRouter(router)
	locationHdl.EnrichRouter(router)

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":"+port, handler))
}
