package utils

import (
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
}

func (a *AppError) ErrorCode() int {
	return a.Code
}

func (a *AppError) Error() string {
	return a.Message
}

func CjsonError(c *gin.Context, err error) {
	c.JSON(err.(*AppError).ErrorCode(), gin.H{"error": err.(*AppError).Error()})
}
