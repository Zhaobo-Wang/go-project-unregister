package controllers

import (
	"github.com/Zhaobo-Wang/go-project-unregister/models"
	"github.com/Zhaobo-Wang/go-project-unregister/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	apiVersion := c.Query("api_version")
	var response services.LLMResponse
	var err error

	switch apiVersion {
	case "v2":
		response, err = services.SendMessageToLLMV2(message)
	default:
		response, err = services.SendMessageToLLMV1(message)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
