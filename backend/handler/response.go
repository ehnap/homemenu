package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Data:    data,
		Message: "success",
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Code:    status,
		Data:    nil,
		Message: message,
	})
}

func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	if id, ok := userID.(int64); ok {
		return id
	}
	return 0
}
