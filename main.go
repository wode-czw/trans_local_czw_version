package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/wodeczw/synk_demo/server_files"
)

func main() {
	//设置并启动gin服务器
	go func() {
		server_files.Start_gin()
	}()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT)

	//启动chrome
	port := "27149"
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+port+"/static/index.html")
	cmd.Start()

	//终端信号关闭应用
	<-chSignal //x := <-chSignal		我不关心读出来的值，会堵塞在这里。
	cmd.Process.Kill()

}
