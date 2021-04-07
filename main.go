package main

import (
	"fmt"
	_ "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"html/template"
	"net/http"
	"os"

	"go-heroku-server/api"
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

//Go application entrypoint
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		//log.Fatal("$PORT must be set")
		port = "5000"
	}

	config.InitializeDatabase()
	config.InitializeRedis()

	config.DBConnection.AutoMigrate(&user.User{}, &types.Address{}, &user.UserNumber{}, &user.Todo{}, &api.File{}, &api.Location{}, &api.LocationImage{}, &location.RestaurantLocation{})

	myRouter := mux.NewRouter()

	myRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	myRouter.HandleFunc("/", hello)

	user.UserEnrichRouter(myRouter)

	myRouter.HandleFunc("/getAllStudents", user.GetAllStudents).Methods("GET")
	myRouter.HandleFunc("/getStudentTodos/{id}", user.GetStudentTodos).Methods("GET")
	myRouter.HandleFunc("/addStudentTodo/{id}", user.AddStudentTodo).Methods("POST")
	myRouter.HandleFunc("/uploadFile", api.UploadFile).Methods("POST")
	myRouter.HandleFunc("/downloadFile/{id}", api.ServeFile).Methods("GET")
	myRouter.HandleFunc("/fileList", api.GetFileList).Methods("GET")
	myRouter.HandleFunc("/getLocations", api.GetLocations).Methods("GET")
	myRouter.HandleFunc("/getLocationImage/{id}", api.GetLocationImage).Methods("GET")
	myRouter.HandleFunc("/saveLocation", api.AddLocation).Methods("POST")
	myRouter.HandleFunc("/getRestaurants", location.GetRestaurantLocations).Methods("GET")
	myRouter.HandleFunc("/getRestaurantByName", location.GetRestaurantByName).Methods("POST")

	fmt.Println("Listening")

	handler := cors.AllowAll().Handler(myRouter)

	fmt.Println(http.ListenAndServe(":"+port, handler))
}
