package route

import (
	"main/controller"
	"net/http"
)

func RegisterVideoRoutes() {
	http.HandleFunc("/videos/create", controller.CreateVideo)
	http.HandleFunc("/videos/get", controller.GetVideo)
}
