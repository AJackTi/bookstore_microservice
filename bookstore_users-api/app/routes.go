package app

import (
	"github.com/AJackTi/bookstore_users-api/controllers/ping"
	"github.com/AJackTi/bookstore_users-api/controllers/users"
)

func route() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.GetUser)
	router.GET("/users/search", users.FindUser)
	router.POST("/users", users.CreateUser)
}
