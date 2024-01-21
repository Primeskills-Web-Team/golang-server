package primeskillsserver

import "github.com/gin-gonic/gin"

type (
	MethodeNotFoundHandler func(c *gin.Context)
	NoRouteHandler         func(c *gin.Context)
	ExceptionHandler       func(c *gin.Context, err any)
	MiddlewareHandler      func(c *gin.Context)
	RouterHandler          func(c *gin.Context)
)
