package server_files

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/wode-czw/tran_tools_czw/config"
	"github.com/wode-czw/tran_tools_czw/server_files/ws"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Start_gin() {
	hub := ws.NewHub()
	go hub.Run()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.POST("/api/v1/files", FilesController)
	router.GET("/api/v1/qrcodes", QrcodesController)
	router.StaticFS("/static", http.FS(staticFiles))
	router.GET("/uploads/:path", UploadsController)
	router.POST("api/v1/texts", TextsController)
	router.GET("api/v1/addresses", AddressesController)
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	router.Run(":" + config.Get_Port())
}
