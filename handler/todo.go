package handler

import (
	"mytodoApp/database/dbHelper"
	"mytodoApp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {

	var req models.CreateTodo

	//	parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalids input",
		})
		return
	}

	//get token from header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
	}
	//validate session to get userID
	userID, err := dbHelper.ValidateSession(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session",
		})
		return
	}
	//userID := userUUId.String()

	// validate date
	if req.ExpiryAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date expiry time over",
		})
		return
	}
	//	create todo
	todo, err := dbHelper.CreateTodo(userID, req.Name, req.Description, req.ExpiryAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create todo",
		})
		return
	}
	//response

	c.JSON(http.StatusCreated, todo)

}

func GetTodos(c *gin.Context) {
	// get token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
		return
	}

	//	validate session
	userID, err := dbHelper.ValidateSession(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session",
		})
		return
	}

	//	get query params
	status := c.Query("status")
	//expiryAt := c.Query("expiryAt")
	search := c.Query("search")

	var completeFilter *bool
	var pendingFilter bool

	switch status {
	case "completed":
		val := true
		completeFilter = &val
	case "incomplete":
		val := false
		completeFilter = &val
	case "pending":
		pendingFilter = true
	}

	//	call database
	todos, err := dbHelper.GetTodos(userID, search, completeFilter, pendingFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	//	response
	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})

}

func GetTodoByID(c *gin.Context) {
	//	get token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
		return
	}
	//	validate session
	userID, err := dbHelper.ValidateSession(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session",
		})
		return
	}
	//	get todo id from query
	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "todo not found",
		})
		return
	}
	//fetch todo from database
	todo, err := dbHelper.GetTodoByID(todoID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "todo not found",
		})
	}
	//response
	c.JSON(http.StatusOK, todo)

}

func UpdateTodoByID(c *gin.Context) {
	//get token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
		return
	}

	//validate session
	userId, err := dbHelper.ValidateSession(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session",
		})
		return
	}

	//get todoid from params
	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "todo id required",
		})
		return
	}

	//parse request
	var req models.UpdateTodo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	//validate date
	if req.ExpiryAt != nil && req.ExpiryAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date expiry time over",
		})
		return
	}

	//update db
	err = dbHelper.UpdateTodo(req, todoID, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//response
	c.JSON(http.StatusOK, gin.H{
		"message": "todo updated successfully",
	})
}

func DeleteTodoByID(c *gin.Context) {

	//get token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token not found",
		})
	}

	//validation session
	userID, err := dbHelper.ValidateSession(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "session invalid",
		})
	}

	//get todo by id
	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "todo id requires",
		})
		return
	}

	//delete
	err = dbHelper.DeleteTodoByID(todoID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	//response
	c.JSON(http.StatusOK, gin.H{
		"message": "todo deleted successfully",
	})
}
