package aria2

import (
	"IxDShell/util"
	"context"
	"fmt"
	rpcc "github.com/zyxar/argo/rpc"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var (
	aria2 rpcc.Client
)

func StartAria2() error {
	//检查目录
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		return err
	}
	//就两层，不做递归了
	distDir := usr.HomeDir + "/Downloads/IxD"
	err = checkAndMkdir(usr.HomeDir + "/Downloads")
	if err != nil {
		return err
	}
	err = checkAndMkdir(distDir)
	if err != nil {
		return err
	}

	//获取一个没有被使用的接口
	freePort, err := util.GetFreePort()
	if err != nil {
		return err
	}
	//启动aria
	err = launchAria2cDaemon(freePort, distDir)
	if err != nil {
		return err
	}
	rpcURI := fmt.Sprintf("http://localhost:%d/jsonrpc", freePort)
	aria2, err = rpcc.New(context.Background(), rpcURI, "", time.Second, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	//defer aria2.Close()
	return nil
}

func launchAria2cDaemon(port int, dist string) (err error) {
	aria2cFile := "aria2c"
	switch runtime.GOOS {
	case "darwin":
		break
	case "windows":
		aria2cFile = "aria2c.exe"
		break
	case "linux":
		break
	}
	//启动aria2的后台服务
	cmdStr := fmt.Sprintf("./%s -d %s --enable-rpc --rpc-listen-all --rpc-listen-port=%d", aria2cFile, dist, port)
	list := strings.Split(cmdStr, " ")
	cmd := exec.Command(list[0], list[1:]...)
	if err = cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	return cmd.Process.Release()
}

func checkAndMkdir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func DownloadUrl(url, name string) (string, error) {
	log.Println(name)
	//增加普通下载
	gid, err := aria2.AddURI(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return gid, nil
}

func DownloadState(gid string) (rpcc.StatusInfo, error) {
	info, err := aria2.TellStatus(gid)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	return info, nil
}
