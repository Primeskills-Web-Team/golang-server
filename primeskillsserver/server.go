package primeskillsserver

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

var logFormatter = &logrus.TextFormatter{
	FullTimestamp: true,
}

func logFormatterSetting() {
	logrus.SetFormatter(logFormatter)
	logrus.SetOutput(os.Stdout)
}

type PrimeskillsServer struct {
	MethodeNotFoundHandler MethodeNotFoundHandler
	NoRouteHandler         NoRouteHandler
	ExceptionHandler       ExceptionHandler
	RouterHandler          RouterHandler
	MiddlewareHandler      MiddlewareHandler
	EngineServer           *gin.Engine
}

func NewPrimeskillsServer() *PrimeskillsServer {
	logFormatterSetting()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	return &PrimeskillsServer{
		EngineServer: r,
	}
}

func (p *PrimeskillsServer) SetStatusMethodNotAllowed(statusMethodNotAllowed MethodeNotFoundHandler) {
	p.MethodeNotFoundHandler = statusMethodNotAllowed
}

func (p *PrimeskillsServer) SetStatusNotFound(notFound NoRouteHandler) {
	p.NoRouteHandler = notFound
}

func (p *PrimeskillsServer) SetException(exceptions ExceptionHandler) {
	p.ExceptionHandler = exceptions
}

func (p *PrimeskillsServer) SetRouters(router RouterHandler) {
	p.RouterHandler = router
}

func (p *PrimeskillsServer) SetMiddleware(middleware MiddlewareHandler) {
	p.MiddlewareHandler = middleware
}

func (p *PrimeskillsServer) RunServer(port string) {
	logrus.Infoln("You can running in", "http://localhost:"+port)

	if p.MethodeNotFoundHandler != nil {
		p.EngineServer.NoMethod(func(context *gin.Context) {
			p.MethodeNotFoundHandler(context)
		})
	} else {
		p.EngineServer.NoMethod(func(context *gin.Context) {
			context.JSON(http.StatusMethodNotAllowed, map[string]interface{}{"error": "Method Not Allowed"})
		})
	}

	if p.NoRouteHandler != nil {
		p.EngineServer.NoRoute(func(context *gin.Context) {
			p.NoRouteHandler(context)
		})
	} else {
		p.EngineServer.NoRoute(func(context *gin.Context) {
			context.JSON(http.StatusNotFound, map[string]interface{}{"error": "Route not found"})
		})
	}

	if p.ExceptionHandler != nil {
		p.EngineServer.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
			p.ExceptionHandler(c, err)
			return
		}))
	} else {
		p.EngineServer.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
			c.JSON(http.StatusNotFound, map[string]interface{}{"error": fmt.Sprintf("%v", err)})
			return
		}))
	}

	err := p.EngineServer.Run(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalln("Failed start server", err.Error())
	}
}
