package routes

import (
	"project/restapi/middelwares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/", helloWorld)

	server.POST("/signup", singup)

	server.POST("/login", login)

	authenticatedRoutes := server.Group("/")
	authenticatedRoutes.Use(middelwares.Authenticate)

	authenticatedRoutes.GET("/events", getAllEventes)

	authenticatedRoutes.GET("/events/:id", getEventByID)

	authenticatedRoutes.POST("/events", createEvent)

	authenticatedRoutes.PUT("/events/:id", updateEvent)

	authenticatedRoutes.DELETE("/events/:id", deleteEvent)
}
