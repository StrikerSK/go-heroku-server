package main

import (
	"fmt"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go-heroku-server/api/files"
	"go-heroku-server/api/location/image"
	"go-heroku-server/api/location/restaurant"
	todoHandlers "go-heroku-server/api/todo/handler"
	todoRepositories "go-heroku-server/api/todo/repository"
	todoServices "go-heroku-server/api/todo/service"
	"go-heroku-server/api/types"
	userAuth "go-heroku-server/api/user/auth"
	userHandlers "go-heroku-server/api/user/handler"
	userRepositories "go-heroku-server/api/user/repository"
	userServices "go-heroku-server/api/user/service"
	"html/template"
	"net/http"
	"os"

	"go-heroku-server/api/location"
	"go-heroku-server/config"
)

func serveMainPage(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("static/index.html"))
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func init() {
	config.GetDatabaseInstance().AutoMigrate(&files.File{}, &types.Address{}, &location.UserLocation{}, &image.LocationImage{}, &restaurant.RestaurantLocation{})
	config.GetCacheInstance()
}

//Go application entrypoint
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "5000"
	}

	userRepository := userRepositories.NewUserRepository(config.GetDatabaseInstance())
	userService := userServices.NewUserService(userRepository)

	userTokenService := userAuth.NewTokenService()
	userMiddleware := userHandlers.NewUserAuthMiddleware(userTokenService)
	userHdl := userHandlers.NewUserHandler(userService, userMiddleware, userTokenService)

	todoRepo := todoRepositories.NewTodoRepository(config.GetDatabaseInstance())
	todoService := todoServices.NewTodoService(todoRepo)
	todoHdl := todoHandlers.NewTodoHandler(userMiddleware, todoService)

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", serveMainPage)

	userHdl.EnrichRouter(router)
	todoHdl.EnrichRouter(router)
	files.EnrichRouteWithFile(router)
	location.EnrichRouteWithLocation(router)

	handler := cors.AllowAll().Handler(router)

	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":"+port, handler))
}
