package shell

import (
	. "IxDShell/config"
	"github.com/zserge/webview"
	"log"
)

func StartWindows() {
	//启动webview
	w := webview.New(webview.Settings{
		Width:                  1100,
		Height:                 618,
		Title:                  "IxD",
		Resizable:              true,
		URL:                    CONF.WebAddr + "/login",
		ExternalInvokeCallback: winHandleRPC,
	})
	defer w.Exit()
	w.Run()
}

func winHandleRPC(w webview.WebView, data string) {
	switch {
	case data == "close":
		w.Terminate()
	case data == "open":
		file := w.Dialog(webview.DialogTypeOpen, 0, "Open file", "")
		log.Println(file)
	case data == "opendir":
		dir := w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", "")
		log.Println(dir)
	}
}
