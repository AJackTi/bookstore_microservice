package app

import (
	"github.com/AJackTi/bookstore_users-api/controllers/ping"
	"github.com/AJackTi/bookstore_users-api/controllers/users"
)

func route() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.Get)
	router.POST("/users", users.Create)
	router.PUT("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.GET("/users/search", users.Search)
}
