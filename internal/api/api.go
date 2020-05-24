package api

import (
	"strings"

	"github.com/biensupernice/krane/internal/api/handler"
	"github.com/biensupernice/krane/internal/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Config : server config
type Config struct {
	RestPort string // port to use for rest api
	LogLevel string // release | debug
}

// Start : api server
func Start(cnf Config) {
	gin.SetMode(strings.ToLower(cnf.LogLevel))

	client := gin.New()

	// Middleware
	client.Use(gin.Recovery())
	client.Use(gin.Logger())
	client.Use(cors.Default())

	// Routes
	client.POST("/health", handler.HealthHandler)
	client.POST("/auth", handler.Auth)
	client.GET("/login", handler.Login)

	client.GET("/sessions", middleware.AuthSessionMiddleware(), handler.GetSessions)

	client.GET("/deployments", middleware.AuthSessionMiddleware(), handler.GetDeployments)
	client.GET("/deployments/:name", middleware.AuthSessionMiddleware(), handler.GetDeployment)
	client.POST("/deployments", middleware.AuthSessionMiddleware(), handler.CreateDeployment)
	client.POST("/deployments/:name/run", middleware.AuthSessionMiddleware(), handler.RunDeployment) // ex. /deployments/:name/run?tag=latest
	client.DELETE("/deployments/:name", middleware.AuthSessionMiddleware(), handler.DeleteDeployment)

	client.GET("/containers", middleware.AuthSessionMiddleware(), handler.ListContainers)
	client.GET("/containers/:containerID/events", handler.ContainerEvents)
	client.PUT("/containers/:containerID/stop", middleware.AuthSessionMiddleware(), handler.StopContainer)
	client.PUT("/containers/:containerID/start", middleware.AuthSessionMiddleware(), handler.StartContainer)

	// --  Websockets -- //

	client.Run(":" + cnf.RestPort)
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func wshandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := wsupgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Failed to set websocket upgrade: %+v", err)
// 		return
// 	}

// 	for {
// 		t, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		conn.WriteMessage(t, msg)
// 	}
// }
