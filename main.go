package main

import (
	"IxDShell/server/aria2"
	"IxDShell/shell"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {

	//单独协程运行启动aria2
	go func() {
		err := aria2.StartAria2()
		if err != nil {
			os.Exit(2)
		}
	}()

	time.AfterFunc(3*time.Second, func() {
		gid, err := aria2.DownloadUrl("http://storage.quickex.com.cn/Fgj6-s-0Z9bWP-AMRNJFK8SXfdmU", "")
		if err != nil {
			return
		}
		log.Println(gid)
	})
	//select {}

	//mac客户端
	shell.StartMac()
	//windows客户端
	//shell.StartWindows()
}
