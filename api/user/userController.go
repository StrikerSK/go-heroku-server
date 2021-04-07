package user

import "github.com/gorilla/mux"

func UserEnrichRouter(router *mux.Router) {
	router.HandleFunc("/getUsers", GetUserList).Methods("GET")
	router.HandleFunc("/user", GetUserDetail).Methods("GET")
	router.HandleFunc("/user", RegisterNewUser).Methods("POST")
	router.HandleFunc("/user", EditUser).Methods("PUT")
	router.HandleFunc("/user/login", Login).Methods("POST")
	router.HandleFunc("/user/welcome", ValidateRequest).Methods("GET")
	//router.HandleFunc("/user/login", LoginUser).Methods("POST")
}
