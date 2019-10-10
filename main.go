package main

import (
	"IxDShell/shell"
	"log"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	//单独协程运行启动aria2
	/*go func() {
		err := aria2.StartAria2()
		if err != nil {
			os.Exit(2)
		}
	}()*/

	/*time.AfterFunc(3*time.Second, func() {
		gid, err := aria2.DownloadUrl("http://storage.quickex.com.cn/lh3AglxbXpUIZR-O5s13UJm3Psei", "")
		if err != nil {
			return
		}
		log.Println(gid)
	})
	select {}*/

	//mac客户端
	shell.StartMac()
	//windows客户端
	//shell.StartWindows()
}
