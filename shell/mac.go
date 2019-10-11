package shell

import (
	. "IxDShell/config"
	"IxDShell/param/upload"
	"IxDShell/service"
	"encoding/json"
	"fmt"
	"github.com/zserge/webview"
	"log"
	"strings"
	"time"
)

func StartMac() {
	//启动webview
	w := webview.New(webview.Settings{
		Width:                  1100,
		Height:                 618,
		Title:                  "IxD",
		Resizable:              true,
		URL:                    fmt.Sprintf("%s/login?t=%d&client=true", CONF.WebAddr, time.Now().Unix()),
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
	case p.Key == "openFileDialog":
		filename := w.Dialog(webview.DialogTypeOpen, 0, "上传单文件", "")
		if filename != "" {
			s := fmt.Sprintf(`externalInvokeOpen("%s")`, filename)
			err := w.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	case p.Key == "openDirDialog":
		dir := w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "上传文件夹", "")
		if dir != "" {
			s := fmt.Sprintf(`externalInvokeOpenDir("%s")`, dir)
			err := w.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	case p.Key == "getLoadingList":
		pathArr, err := service.LocalUpLoadingList()
		if err != nil {
			return
		}
		if len(pathArr) > 0 {
			s := fmt.Sprintf(`externalInvokeLoadingList("%s")`, strings.Join(pathArr, ","))
			err := w.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	case p.Key == "getLoadingProgress":
		mapStr, err := service.LocalUpLoadingProgress()
		if err != nil {
			return
		}
		if mapStr != "" {
			s := fmt.Sprintf(`externalInvokeLoadingProgress(%s)`, mapStr)
			err := w.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	//上传单个文件
	case p.Key == "clientUploadOne":
		//拆解参数
		paramArr := strings.Split(p.Value, "||")
		param := new(upload.One)
		param.Pid = paramArr[0]
		param.LocalPath = paramArr[1]
		err := service.UploadOne(param, paramArr[2])
		res := 0
		if err != nil {
			res = 1
		}
		s := fmt.Sprintf(`externalInvokeClientUploadOne(%d)`, res)
		err = w.Eval(s)
		if err != nil {
			log.Println(err)
		}
	case p.Key == "restartTask":
		_, err := service.UploadRestartTask(p.Value)
		if err != nil {
			return
		}
	case p.Key == "downloadFile":
		paramArr := strings.Split(p.Value, "||")
		err := service.DownloadFile(paramArr[0], paramArr[1])
		if err != nil {
			return
		}
	}
}
