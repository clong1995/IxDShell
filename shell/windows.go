package shell

import (
	. "IxDShell/config"
	"IxDShell/param/upload"
	"IxDShell/service"
	"fmt"
	"github.com/gen2brain/dlgs"
	"github.com/zserge/lorca"
	"log"
	"strings"
	"time"
)

func StartWindows() {
	localhost := fmt.Sprintf("%s/login?t=%d&client=true", CONF.WebAddr, time.Now().Unix())
	ui, err := lorca.New(localhost, "", 1100, 618, "--autoplay-policy=no-user-gesture-required")
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	//打开单个文件
	err = ui.Bind("openFileDialog", func() {
		//filename, ret, err := dlgs.File("Choose directory", "", true)
		filename, ret, err := dlgs.File("Select file", "", false)
		if err != nil {
			log.Println(err)
		}
		if filename != "" && ret {
			filename = strings.Replace(filename, "\\", "\\\\", -1)
			s := fmt.Sprintf(`externalInvokeOpen("%s")`, filename)
			_ = ui.Eval(s)
		}
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	//上传单个文件
	err = ui.Bind("clientUploadOne", func(pid, localPath, Authorization string) {
		p := new(upload.One)
		p.Pid = pid
		p.LocalPath = localPath
		err := service.UploadOne(p, Authorization)
		res := 0
		if err != nil {
			res = 1
		}
		s := fmt.Sprintf(`externalInvokeClientUploadOne(%d)`, res)
		_ = ui.Eval(s)
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	//获取正在上传的文件
	err = ui.Bind("getLoadingList", func() {
		pathArr, err := service.LocalUpLoadingList()
		if err != nil {
			return
		}
		if len(pathArr) > 0 {
			s := fmt.Sprintf(`externalInvokeLoadingList("%s")`, strings.Join(pathArr, ","))
			_ = ui.Eval(s)
		}
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	//获取正在上传的文件进度
	err = ui.Bind("getLoadingProgress", func() {
		mapStr, err := service.LocalUpLoadingProgress()
		if err != nil {
			return
		}
		if mapStr != "" {
			s := fmt.Sprintf(`externalInvokeLoadingProgress(%s)`, mapStr)
			_ = ui.Eval(s)
		}
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	//重启上传
	err = ui.Bind("restartTask", func(Authorization string) {
		_, err := service.UploadRestartTask(Authorization)
		if err != nil {
			return
		}
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	//下载
	err = ui.Bind("downloadFile", func(etag, name string) {
		err := service.DownloadFile(etag, name)
		if err != nil {
			return
		}
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	//获取下载状态
	err = ui.Bind("getDownloadProgress", func(Authorization string) {

		service.DownloadProgress(Authorization)

		/*if err != nil {
			return
		}
		if mapStr != "" {
			s := fmt.Sprintf(`externalInvokeDownloadProgress(%s)`, mapStr)
			_ = ui.Eval(s)
		}*/
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	<-ui.Done()
}
