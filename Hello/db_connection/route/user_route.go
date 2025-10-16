package route

import (
	"main/controller"
	"net/http"
)

func RegisterUserRoutes() {
	http.HandleFunc("/users", controller.CreateUser)
	http.HandleFunc("/users/get", controller.GetUser)
	http.HandleFunc("/users/delete", controller.DeleteUser)
	http.HandleFunc("/users/update", controller.UpdateUser)
	http.HandleFunc("/users/login", controller.LoginUser)
}
