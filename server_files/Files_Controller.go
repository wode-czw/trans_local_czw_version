package server_files

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw")
	if err != nil {
		log.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	if err != nil {
		log.Fatal(err)
	}
	// 获取当前时间并格式化
	currentDate := time.Now().Format("20060102")
	currentTime := time.Now().Format("150405")

	// 使用下划线将日期和时间分开
	formattedTime := currentDate + "_" + currentTime

	// 构建新的文件名：时间+原文件名+扩展名
	filename := fmt.Sprintf("%s___%s", formattedTime, file.Filename)

	uploads := filepath.Join(dir, "uploads")
	err = os.MkdirAll(uploads, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fullpath := path.Join("uploads", filename)
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	log.Printf("url" + "/" + fullpath)
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) //保存好后会返回ok的状态码以及地址

}
