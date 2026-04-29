package handler

import (
	"mytodoApp/database/dbHelper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTodos(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")

	todos, err := dbHelper.GetAllTodos(search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

func GetAllUsers(c *gin.Context) {

	search := c.Query("search")

	users, err := dbHelper.GetAllUsers(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func SuspendUser(c *gin.Context) {

	//	get user id from url
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id is required",
		})
		return
	}

	//call db handler
	err := dbHelper.SuspendUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user suspended successfully",
	})
}

func UnsuspendUser(c *gin.Context) {

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id is required",
		})
		return
	}

	err := dbHelper.UnsuspendUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user unsuspended successfully",
	})
}
