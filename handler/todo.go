package handler

import (
	"mytodoApp/database/dbHelper"
	"mytodoApp/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {

	var req models.CreateTodo

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid inputs",
		})
		return
	}

	userID := c.GetString("userID")

	//handle error
	//if req.Name == "" || req.Description == "" || req.ExpiryAt.IsZero() {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"#error": "all these fields are required ",
	//	})
	//	return
	//}

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
	c.JSON(http.StatusCreated, todo)
}

func GetTodos(c *gin.Context) {

	userID := c.GetString("userID")

	status := c.Query("status")
	search := c.Query("search")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	if status != "" && status != "completed" && status != "incomplete" && status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status query",
		})
	}

	todos, err := dbHelper.GetTodos(userID, search, status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"page":  page,
		"limit": limit,
	})
}

func GetTodoByID(c *gin.Context) {

	userID := c.GetString("userID")

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "todo not found",
		})
		return
	}

	todo, err := dbHelper.GetTodoByID(todoID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "todo not found",
		})
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodoByID(c *gin.Context) {

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "todo id required",
		})
		return
	}

	var req models.UpdateTodo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	if req.ExpiryAt != nil && req.ExpiryAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date expiry time over",
		})
		return
	}
	userID := c.GetString("userID")

	//update db
	err := dbHelper.UpdateTodo(req, todoID, userID)
	if err != nil {
		//fmt.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "todo updated successfully",
	})
}

func DeleteTodoByID(c *gin.Context) {

	userID := c.GetString("userID")

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "todo id requires",
		})
		return
	}

	err := dbHelper.DeleteTodoByID(todoID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "todo deleted successfully",
	})
}
