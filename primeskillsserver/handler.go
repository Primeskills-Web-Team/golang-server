package primeskillsserver

import "github.com/gin-gonic/gin"

type (
	MethodeNotFoundHandler func(c *gin.Context)
	NoRouteHandler         func(c *gin.Context)
	ExceptionHandler       func(c *gin.Context, err interface{}, resource string)
	MiddlewareHandler      func(c *gin.Context)
	RouterHandler          func(c *gin.Context)
)
