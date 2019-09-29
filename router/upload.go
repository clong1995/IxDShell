package router

import (
	"IxDShell/handler"
	"net/http"
)

func init() {
	//上传
	http.HandleFunc("/upload/one", handler.UploadOne)
	//重新载入任务列表
	//http.HandleFunc("/upload/restartTask", handler.UploadRestartTask)
}
