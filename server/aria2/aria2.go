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
	"path/filepath"
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
	currDir, err := util.CurrDir()
	if err != nil {
		log.Println(err)
		return
	}
	if strings.HasPrefix(currDir, "/private") {
		currDir = "/"
	}
	aria2cFilePath := ""
	switch runtime.GOOS {
	case "darwin":
		aria2cFilePath = filepath.Join(currDir, "aria2c")
		break
	case "windows":
		aria2cFilePath = filepath.Join(currDir, "bin", "aria2c.exe")
		break
	case "linux":
		break
	}

	//启动aria2的后台服务
	cmdStr := fmt.Sprintf(".%s -d %s --enable-rpc --rpc-listen-all --rpc-listen-port=%d", aria2cFilePath, dist, port)
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
	//增加普通下载
	gid, err := aria2.AddURI(url)
	if err != nil {
		log.Println(err)
		return "", err
	}

	//启动线程检查下载状态，用来重命名
	go func() {
		for range time.Tick(1 * time.Second) {
			msg, err := aria2.TellStatus(gid)
			if err != nil {
				fmt.Println(err)
				break
			}
			if msg.Status == "complete" {
				//下载完成，重命名
				err := os.Rename(msg.Files[0].Path, filepath.Join(msg.Dir, name))
				if err != nil {
					log.Println(err)
				}
				break
			} else {
				//log.Println(msg)
			}
		}
	}()
	return gid, nil
}

func DownloadRestart() {
	//获取所有session
	//启动下载
}

/*func DownloadProgress(gid string) (rpcc.StatusInfo, error) {
	info, err := aria2.TellStatus(gid)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	return info, nil
}*/

func DownloadProgress() {
	//aria2会自动断点下载，可以根据未完成的文件的名字（.aria2）判断是否下发，没下完的名字是etag。
	//然后内存记录gid，查询进度，下载完成，删除内存的gid
	//aria2没法关闭。。。

	//查询下载进度 163abb8aed1dd16a
	log.Println("查询下载进度")
	//info, err := aria2.GetGlobalStat()

	info, err := aria2.TellStatus("163abb8aed1dd16a")
	if err != nil {
		log.Println(err)
	}
	log.Println(info)

}
