package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ginEngine *gin.Engine
}

func New() *Server {
	engine := gin.New()

	// Pasang middleware bawaan Gin dan Custom Error Handler buatan kita
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.Use(CustomErrorHandler())

	return &Server{ginEngine: engine}
}

// Start menghidupkan server
func (s *Server) Start(addr string) error {
	return s.ginEngine.Run(addr)
}

// Use untuk memasang middleware global (seperti CORS)
func (s *Server) Use(middlewares ...gin.HandlerFunc) {
	s.ginEngine.Use(middlewares...)
}

// HealthHandler untuk mengecek status nyala/mati server
func (s *Server) HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// GET, POST, DELETE meneruskan perintah ke framework Gin
// (Kita pisahkan Handler utama 'h' dengan array Middleware 'mw')
func (s *Server) GET(path string, h gin.HandlerFunc, mw ...gin.HandlerFunc) {
	s.ginEngine.GET(path, append(mw, h)...)
}
func (s *Server) POST(path string, h gin.HandlerFunc, mw ...gin.HandlerFunc) {
	s.ginEngine.POST(path, append(mw, h)...)
}
func (s *Server) DELETE(path string, h gin.HandlerFunc, mw ...gin.HandlerFunc) {
	s.ginEngine.DELETE(path, append(mw, h)...)
}

// Group membungkus fitur grup URL
func (s *Server) Group(prefix string, mw ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{group: s.ginEngine.Group(prefix, mw...)}
}

// RouterGroup struct untuk meniru *gin.RouterGroup
type RouterGroup struct {
	group *gin.RouterGroup
}

func (rg *RouterGroup) GET(path string, h gin.HandlerFunc, mw ...gin.HandlerFunc) {
	rg.group.GET(path, append(mw, h)...)
}
func (rg *RouterGroup) POST(path string, h gin.HandlerFunc, mw ...gin.HandlerFunc) {
	rg.group.POST(path, append(mw, h)...)
}
