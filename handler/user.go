package handler

import (
	"mytodoApp/database/dbHelper"
	"mytodoApp/models"
	"mytodoApp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user models.RegisterUser

	//   parse request body gin read incoming json and map it into struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	// check if user exists or not
	exist, err := dbHelper.IsUserExist(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error",
		})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "user already exist",
		})
		return
	}

	//create user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	_, err = dbHelper.CreateUser(user.Name, user.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//get response
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}

// login
func LoginUser(c *gin.Context) {
	var req models.LoginUser

	//   parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}
	//get user from DB
	userID, hashedPassword, err := dbHelper.GetUserBYEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	//check password
	if err := utils.CheckPasswordHash(hashedPassword, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
	}

	//	create session
	sessionID, err := dbHelper.CreateUserSession(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}

	//    response
	c.JSON(http.StatusOK, gin.H{
		"token":   sessionID,
		"message": "login successfully",
	})
}

// logout
func LogoutUser(c *gin.Context) {
	token := c.GetString("token")
	//if token == "" {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"error": "unauthorized missing token",
	//	})
	//	return
	//}
	err := dbHelper.ArchiveUserSession(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid session token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user logout successfully",
	})
}
