package routers

import (
	"resapi/internal/app/handlers"
	"resapi/internal/domain/infs/database"
	"resapi/internal/domain/repo"

	"github.com/gin-gonic/gin"
)

func PubRoutSetup() *gin.Engine {

	userRepo := repo.NewUserRepo(database.DB)

	userHandler := handlers.NewUserHandler(*userRepo)

	r := gin.Default()

	r.LoadHTMLGlob("template/*")

	//Public routers

	r.GET("/", handlers.IndexPageHandler)

	r.POST("/reg", userHandler.CreateUserHandler) //Create user
	r.GET("/user/:id", userHandler.GetUser)       //Get user
	r.PUT("/user/:id", userHandler.UpdateUserHandler)
	r.DELETE("/user/:id", userHandler.DeleteUserHandler)
	return r
}
