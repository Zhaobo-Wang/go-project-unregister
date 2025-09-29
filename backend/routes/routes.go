package routes

import (
	"github.com/Zhaobo-Wang/go-projects/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	
	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)
	
	api.GET("/todos", controllers.GetTodos)
	api.POST("/todos", controllers.CreateTodo)
	api.GET("/todos/:id", controllers.GetTodo)
	api.PATCH("/todos/:id", controllers.UpdateTodo)
	api.PUT("/todos/:id", controllers.UpdateTodo)
	api.DELETE("/todos/:id", controllers.DeleteTodo)

	api.GET("/user-profile", controllers.GetUser)
	
	return r
}
