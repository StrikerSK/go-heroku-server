package main

import (
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go-heroku-server/api/files"
	"go-heroku-server/api/location/image"
	"go-heroku-server/api/location/restaurant"
	"go-heroku-server/api/todo"
	"go-heroku-server/api/types"
	"go-heroku-server/config"
	"html/template"
	"log"
	"net/http"
	"os"

	"go-heroku-server/api/location"
	"go-heroku-server/api/user"
)

func serveMainPage(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles("static/index.html"))
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	config.GetDatabaseInstance().AutoMigrate(&user.User{}, &files.File{}, &types.Address{}, &location.UserLocation{}, &image.LocationImage{}, &restaurant.RestaurantLocation{})
	config.GetCacheInstance()
	user.InitializeUsers()
}

//Go application entrypoint
func main() {
	if err := createServer().ListenAndServe(); err != nil {
		log.Printf("Server starting: %s\n", err.Error())
		os.Exit(1)
		return
	}
}

func createServer() *http.Server {
	server := &http.Server{
		Addr:    ":" + resolvePort(),
		Handler: createHandlers(),
	}

	log.Print("Server has been prepared")
	return server
}

func resolvePort() string {
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("Using default port 5000, because PORT has not been set!")
		port = "5000"
	}

	log.Printf("Port [%s] has been set for server!\n", port)
	return port
}

func createHandlers() (handler http.Handler) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	myRouter.HandleFunc("/", serveMainPage)

	user.EnrichRouterWithUser(myRouter)
	files.EnrichRouteWithFile(myRouter)
	todo.EnrichRouteWithTodo(myRouter)
	location.EnrichRouteWithLocation(myRouter)

	handler = cors.AllowAll().Handler(myRouter)
	return
}
