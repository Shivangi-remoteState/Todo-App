package server

import (
	"mytodoApp/handler"
	"mytodoApp/middleware"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	//routes
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/logout", handler.LogoutUser)
	//r.GET("/todos", handler.GetTodos)
	todo := auth.Group("/")
	{
		todo.POST("/todo", handler.CreateTodo)
		todo.GET("/todos", handler.GetTodos)
		todo.GET("/todo/:id", handler.GetTodoByID)
		todo.PUT("/todo/:id", handler.UpdateTodoByID)
		todo.DELETE("/todo/:id", handler.DeleteTodoByID)

	}

	admin := auth.Group("/admin")
	admin.Use(middleware.AdminOnly())
	admin.GET("/todos", handler.GetAllTodos)
	admin.GET("/users", handler.GetAllUsers)
	admin.PATCH("/user/:id/suspend", handler.SuspendUser)
	admin.PATCH("/user/:id/unsuspend", handler.UnsuspendUser)

	//r.Run("/todo", handler.CreateTodo)

	if err := r.Run(":8080"); err != nil {
		panic(err.Error())
	}
}
