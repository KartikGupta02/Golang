package route

import (
	"main/controller"
	"main/middleware"
	"net/http"
)

func RegisterUserRoutes() {
	http.HandleFunc("/", controller.GetHtmlData)

	http.HandleFunc("/users", controller.CreateUser)
	http.HandleFunc("/get_all_users", controller.GetUsers)
	http.HandleFunc("/users/get", controller.GetUser)
	http.HandleFunc("/users/delete", middleware.AuthMiddleware(controller.DeleteUser))
	http.HandleFunc("/users/update", middleware.AuthMiddleware(controller.UpdateUser))
	http.HandleFunc("/users/login", controller.LoginUser)
}
