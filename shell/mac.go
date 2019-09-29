package shell

import (
	. "IxDShell/config"
	"encoding/json"
	"fmt"
	"github.com/zserge/webview"
	"log"
)

func StartMac() {
	//启动webview
	w := webview.New(webview.Settings{
		Width:                  1100,
		Height:                 618,
		Title:                  "IxD",
		Resizable:              true,
		URL:                    CONF.WebAddr + "/login?client=true",
		ExternalInvokeCallback: macHandleRPC,
	})
	defer w.Exit()
	w.Run()
}

//参数
type param struct {
	//本地要上传的文件
	Key   string `json:"key"`
	Value string `json:"value"`
}

func macHandleRPC(w webview.WebView, data string) {
	b := []byte(data)
	p := param{}
	err := json.Unmarshal(b, &p)
	if err != nil {
		log.Println(err)
	}
	switch {
	case p.Key == "open":
		file := w.Dialog(webview.DialogTypeOpen, 0, "上传单文件", "")
		if file != "" {
			s := fmt.Sprintf(`externalInvokeOpen("%s")`, file)
			err := w.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	case p.Key == "openDir":
		dir := w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "上传文件夹", "")
		if dir != "" {
			log.Println(dir)
			s := fmt.Sprintf(`externalInvokeOpenDir("%s")`, dir)
			err := w.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
