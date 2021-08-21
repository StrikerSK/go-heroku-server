package main

import (
	"fmt"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go-heroku-server/api/files"
	"go-heroku-server/api/location/image"
	"go-heroku-server/api/location/restaurant"
	"go-heroku-server/api/todo"
	"go-heroku-server/api/types"
	"html/template"
	"net/http"
	"os"

	"go-heroku-server/api/location"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
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
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "5000"
	}

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	myRouter.HandleFunc("/", serveMainPage)

	user.EnrichRouterWithUser(myRouter)
	files.EnrichRouteWithFile(myRouter)
	todo.EnrichRouteWithTodo(myRouter)
	location.EnrichRouteWithLocation(myRouter)

	handler := cors.AllowAll().Handler(myRouter)

	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":"+port, handler))
}
