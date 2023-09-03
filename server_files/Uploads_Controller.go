package server_files

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {

		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe)
		uploads := filepath.Join(dir, "uploads")

		target := filepath.Join(uploads, path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target) //c能够给前端发送任何类型的文件。
	} else {
		c.Status(http.StatusNotFound)
	}
}
