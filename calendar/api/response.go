package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GenericResponse struct {
	Message string `json:"message"`
}

func UnauthorizedA(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, GenericResponse{
		Message: msg,
	})
}

func BadRequestA(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, GenericResponse{
		Message: err.Error(),
	})
}

func BadJSONA(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, GenericResponse{
		Message: "failed to parse request body",
	})
}

func ForbiddenA(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, GenericResponse{
		Message: fmt.Sprintf("%s access denied", msg),
	})
}

func NotFoundA(c *gin.Context, entity string) {
	c.AbortWithStatusJSON(http.StatusNotFound, GenericResponse{
		Message: fmt.Sprintf("%s not found", entity),
	})
}

func ServerErrorA(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, GenericResponse{
		Message: err.Error(),
	})
}
