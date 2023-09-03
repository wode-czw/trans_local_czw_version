package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	//设置并启动gin服务器
	go func() {
		Start_gin()
	}()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT)

	//启动chrome
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+port+"/static/index.html")
	cmd.Start()

	//终端信号关闭应用
	<-chSignal //x := <-chSignal		我不关心读出来的值，会堵塞在这里。
	cmd.Process.Kill()

}

func TextsController(c *gin.Context) {
	var json struct {
		Raw string
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
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

		filename := formattedTime + uuid.New().String()
		uploads := filepath.Join(dir, "uploads")
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")
		err = os.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}

}

func AddressesController(c *gin.Context) {

	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

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

func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256) //这里png是得到的图像的2进制的数据
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

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

func Start_gin() {
	port := "27149"
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.POST("/api/v1/files", FilesController)
	router.GET("/api/v1/qrcodes", QrcodesController)
	router.StaticFS("/static", http.FS(staticFiles))
	router.GET("/uploads/:path", UploadsController)
	router.POST("api/v1/texts", TextsController)
	router.GET("api/v1/addresses", AddressesController)
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

	router.Run(":" + port)
}
