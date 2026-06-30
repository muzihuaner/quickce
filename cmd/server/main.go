package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"quickce/pkg/api"
	"quickce/pkg/config"
	"quickce/pkg/store"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if mode := os.Getenv("GIN_MODE"); mode == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	s := store.New()
	r := gin.Default()

	handler := api.NewHandler(int64(cfg.MaxConcurr), s)
	handler.RegisterRoutes(r)

	agentHandler := api.NewAgentHandler(s)
	agentHandler.RegisterRoutes(r)

	distPath := os.Getenv("STATIC_DIR")
	if distPath == "" {
		distPath = filepath.Join("frontend", "dist")
	}
	if _, err := os.Stat(distPath); err == nil {
		r.Use(serveStaticFallback(distPath))
		log.Printf("Serving frontend from %s", distPath)
	} else {
		r.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "QuickCE API is running. Build the frontend to serve the UI.")
		})
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("QuickCE server starting on %s", addr)
	log.Printf("Max concurrent probes: %d", cfg.MaxConcurr)
	log.Printf("Default timeout: %d seconds", cfg.DefaultTimeout)

	if err := r.Run(addr); err != nil {
		log.Fatalf("server startup failed: %v", err)
	}
}

func serveStaticFallback(root string) gin.HandlerFunc {
	fs := http.FileServer(http.Dir(root))
	return func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			c.Next()
			return
		}
		path := filepath.Join(root, c.Request.URL.Path)
		if _, err := os.Stat(path); err == nil {
			fs.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		c.Request.URL.Path = "/"
		fs.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
