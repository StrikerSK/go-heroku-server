package main

import (
	"fmt"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go-heroku-server/api/files"
	"go-heroku-server/api/todo"
	"html/template"
	"net/http"
	"os"

	"go-heroku-server/api/location"
	"go-heroku-server/api/types"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
)

func hello(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("templates/welcome-page.html"))
	if err := templates.ExecuteTemplate(w, "welcome-page.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func init() {
	config.InitializeDatabase()
	config.InitializeRedis()

	config.DBConnection.AutoMigrate(&user.User{}, &files.File{}, &types.Address{}, &location.Location{}, &location.LocationImage{}, &location.RestaurantLocation{})
	user.InitAdminUser()
	user.InitCommonUser()
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
	myRouter.HandleFunc("/", hello)

	user.EnrichRouterWithUser(myRouter)
	files.EnrichRouteWithFile(myRouter)
	todo.EnrichRouteWithTodo(myRouter)
	location.EnrichRouteWithLocation(myRouter)

	myRouter.HandleFunc("/getRestaurants", location.GetRestaurantLocations).Methods("GET")
	myRouter.HandleFunc("/getRestaurantByName", location.GetRestaurantByName).Methods("POST")

	fmt.Println("Listening")

	handler := cors.AllowAll().Handler(myRouter)

	fmt.Println(http.ListenAndServe(":"+port, handler))
}
