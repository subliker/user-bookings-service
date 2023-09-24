package main

import (
	"github.com/subliker/backendproj/route"

	docs "github.com/subliker/backendproj/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// @title CyberZoneDev test REST API project
// @description This rest api is designed to work with the PostgreSQL database. There are two main entities: User and Booking. One user can have multiple Bookings

var DataBase = route.DataBase

func SetupRouter() *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api"

	router.GET("/api/user/:id", route.GetUserDataById)
	router.POST("/api/user", route.AddNewUser)
	router.DELETE("/api/user/:id", route.DeleteUserDataByID)
	router.PUT("/api/user/:id", route.UpdateUserDataById)

	router.GET("/api/booking/:id", route.GetBookingDataById)
	router.GET("/api/booking", route.GetBookings)
	router.POST("/api/booking", route.AddNewBooking)
	router.DELETE("/api/booking/:id", route.DeleteBookingByID)
	router.PUT("/api/booking/:id", route.UpdateBookingDataById)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func main() {
	router := SetupRouter()
	router.Run(":8000")
}
