package routes

import (
	"github.com/gin-gonic/gin"
	"UserService/internal/handler"
)

type RegisterHandler = handler.RegisterHandler

func Setup() *gin.Engine {
    r := gin.Default()

    auth := r.Group("/api")
    reg := RegisterHandler{}
    auth.POST("/updateUser", reg.Handle)

    return r
}