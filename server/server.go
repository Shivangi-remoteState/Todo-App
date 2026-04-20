package server

import (
	"mytodoApp/handler"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	//routes
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	r.PATCH("/logout", handler.LogoutUser)
	//r.GET("/todos", handler.GetTodos)
	todo := r.Group("/")
	{
		todo.POST("/todo", handler.CreateTodo)
		todo.GET("/todos", handler.GetTodos)
		todo.GET("/todo/:id", handler.GetTodoByID)
		todo.PUT("/todo/:id", handler.UpdateTodoByID)
		todo.DELETE("/todo/:id", handler.DeleteTodoByID)

	}
	//r.POST("/todo", handler.CreateTodo)
	r.Run(":8080")
}
