package router

import (
	"IxDShell/handler"
	"net/http"
)

func init() {
	//TODO 后期改造成全部走js&&go方式
	//上传
	http.HandleFunc("/upload/one", handler.UploadOne)
}
