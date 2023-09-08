package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/wode-czw/tran_tools_czw/config"
	"github.com/wode-czw/tran_tools_czw/server_files"
	//"wodeczw/synk_demo/server_files"
)

func main() {

	go server_files.Start_gin() //设置并启动gin服务器
	go start_browser()          //用locra启动chrome的程序界面

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT)
	//终端信号关闭应用
	<-chSignal //x := <-chSignal		我不关心读出来的值，会堵塞在这里。

}

func start_browser() {
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+config.Get_Port()+"/static/index.html")
	cmd.Start()

}
