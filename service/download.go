package service

import (
	. "IxDShell/config"
	"IxDShell/server/aria2"
	"log"
)

func DownloadFile(etag, name string) error {
	//log.Println(CONF.QiniuAddr + etag)
	//增加下载
	gid, err := aria2.DownloadUrl(CONF.QiniuAddr+etag, name)
	if err != nil {
		return err
	}
	log.Println(gid)
	return nil
}

func DownloadProgress(Authorization string) {
	//查询aria2
	aria2.DownloadProgress()
}
