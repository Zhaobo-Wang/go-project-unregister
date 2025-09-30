package controllers

import (
	"net/http"

	"github.com/Zhaobo-Wang/go-project-unregister/database"
	"github.com/Zhaobo-Wang/go-project-unregister/models"
	"github.com/gin-gonic/gin"
)

// GetUser 获取当前登录用户的信息
func GetUser(c *gin.Context) {
	var user models.User
	if err := database.DB.First(&user, defaultUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}
