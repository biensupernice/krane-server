package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response : response type
type Response struct {
	Success bool        `json:"success"`
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
}

// Ok : response with status code 200
func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Code:    http.StatusOK,
		Data:    data,
	})
	return
}

// Created : response with status code 200
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, &Response{
		Success: true,
		Code:    http.StatusCreated,
		Data:    data,
	})
	return
}

// BadRequest : response with status code 400
func BadRequest(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, &Response{
		Success: false,
		Code:    http.StatusBadRequest,
		Data:    map[string]interface{}{"error": err},
	})
	return
}

// Unauthorized response with status code 400
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, &Response{
		Success: false,
		Code:    http.StatusUnauthorized,
		Data:    map[string]string{"error": "Unauthorized request"},
	})
	return
}
