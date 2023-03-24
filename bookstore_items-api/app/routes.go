package app

import (
	"net/http"

	"github.com/AJackTi/bookstore_items-api/controllers"
)

func routes() {
	router.HandleFunc("/items", controllers.ItemsController.Create).Methods(http.MethodPost)
	router.HandleFunc("/ping", controllers.PingController.Ping).Methods(http.MethodGet)
}
