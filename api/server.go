package api

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"stg-go-websocket-server/routes"
	"stg-go-websocket-server/util"
	"stg-go-websocket-server/ws"
	"strings"
)

type Server struct {
	config  util.Config
	manager *ws.Manager
	router  *gin.Engine
}

func NewServer(config util.Config, m *ws.Manager) (*Server, error) {
	server := &Server{
		config:  config,
		manager: m,
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, ".js") || strings.HasSuffix(c.Request.URL.Path, ".mjs") {
			c.Header("Content-Type", "application/javascript")
		}
	})
	r.Use(static.Serve("/", static.LocalFile("ws-front-end/dist/", true)))
	api := r.Group("/ws")
	api.GET("", s.manager.HandleWS)

	routes.SetupApiRoutes("/api", r, s.manager)
	routes.SetupAdminRoutes("/admin", r, s.manager)
	s.router = r
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start() error {
	return s.router.Run(":" + s.config.PORT)
}
