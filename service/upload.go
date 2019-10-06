package service

import (
	. "IxDShell/common"
	. "IxDShell/config"
	"IxDShell/param/upload"
	"IxDShell/util"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"math"
	"mime"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type progressRecord struct {
	Progresses []storage.BlkputRet `json:"progresses"`
}
type qiniuFileInfoParam struct {
	Etag string
}
type addFileRes struct {
	Etag     string  `json:"etag"`
	Size     float64 `json:"size"`
	MimeType string  `json:"mimeType"`
	Name     string  `json:"name"`
	Pid      string  `json:"pid"`
	State    int     `json:"state"`
	Local    string  `json:"local"`
}

type qiniuFileInfoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data qiniuFileInfoRespData
}

type qiniuFileInfoRespData struct {
	Hash     string  `json:"Hash"`
	Fsize    float64 `json:"FSize"`
	MimeType string  `json:"MimeType"`
}
type uploadFinishRes struct {
	Id string `json:"id"`
}

type getUpKeyResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UpToken string `json:"upToken"`
	}
}
type addFileResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func exists(ckPath string) bool {
	_, err := os.Stat(ckPath)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

//获取上传的key
func getUpKey(Authorization string) (string, error) {
	header := map[string]string{
		"Authorization": Authorization,
	}
	body, err := util.HttpPostJson(CONF.ServerAddr+"/qiniu/key", nil, header)
	if err != nil {
		return "", err
	}

	r := new(getUpKeyResp)
	//解析json
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if r.Code != 0 {
		log.Println(err)
		return "", fmt.Errorf(JSON_STR)
	}

	return r.Data.UpToken, nil
}

//七牛文件信息
func qiniuFileInfo(etag, Authorization string) (interface{}, error) {
	p := qiniuFileInfoParam{
		Etag: etag,
	}
	header := map[string]string{
		"Authorization": Authorization,
	}
	body, err := util.HttpPostJson(CONF.ServerAddr+"/qiniu/fileInfo", p, header)
	if err != nil {
		return nil, err
	}
	r := new(qiniuFileInfoResp)
	//解析json
	err = json.Unmarshal(body, r)
	if r.Code != 0 {
		return nil, fmt.Errorf("no such file or directory")
	}
	return r.Data, nil
}

//上传一个
func UploadOne(p *upload.One, Authorization string) error {
	//检查文件在本地是否存在
	fileInfo, statErr := os.Stat(p.LocalPath)
	if statErr != nil {
		log.Println(statErr)
		return statErr
	}
	//计算文件hash
	etag, err := util.GetEtag(p.LocalPath)
	if err != nil {
		return err
	}
	//认证
	header := map[string]string{
		"Authorization": Authorization,
	}

	//===>查询是否在七牛存在，【触发秒传】
	qiniuFileInfoData, err := qiniuFileInfo(etag, Authorization)
	if err == nil {
		log.Println("触发秒传")
		//直接增加数据库记录
		qfid := qiniuFileInfoData.(qiniuFileInfoRespData)
		afs := addFileRes{
			Etag:     qfid.Hash,
			Name:     fileInfo.Name(),
			MimeType: qfid.MimeType,
			Size:     qfid.Fsize / math.Pow(1024, 2),
			Pid:      p.Pid,
			State:    0,
		}
		_, err := util.HttpPostJson(CONF.ServerAddr+"/file/addFile", afs, header)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	} else {
		log.Println(err)
	}

	//===>否则，开协程去上传，【断点续传】
	//获取上传凭证
	upKey, err := getUpKey(Authorization)
	if err != nil {
		log.Println(err)
		return err
	}
	//更新数据
	mimetype := mime.TypeByExtension(path.Ext(p.LocalPath))
	afs := addFileRes{
		Etag:     etag,
		Name:     fileInfo.Name(),
		Size:     float64(fileInfo.Size()) / math.Pow(1024, 2),
		MimeType: mimetype,
		Pid:      p.Pid,
		State:    2,
		Local:    p.LocalPath,
	}
	body, err := util.HttpPostJson(CONF.ServerAddr+"/file/addFile", afs, header)
	if err != nil {
		log.Println(err)
		return err
	}
	//解析json
	r := new(addFileResp)
	err = json.Unmarshal(body, r)
	if r.Code != 0 {
		return fmt.Errorf(ERR_STR)
	}

	//后台上传
	//文件id，留着上传完了，更新状态用
	go uploadFile(p.LocalPath, upKey, etag, r.Data, Authorization)
	return nil
}

//断点续传
func uploadFile(localFile, upKey, etag, id, Authorization string) error {
	log.Println("触发断点续传！")
	//==进度
	// 必须仔细选择一个能标志上传唯一性的 recordKey 用来记录上传进度
	// 我们这里采用 md5(key+local_path+local_file_last_modified)+".progress" 作为记录上传进度的文件名
	fileInfo, _ := os.Stat(localFile)
	fileSize := fileInfo.Size()
	//recordKey := md5Hex(fmt.Sprintf("%s", etag))
	// 指定的进度文件保存目录，实际情况下，请确保该目录存在，而且只用于记录进度文件
	recordDir := os.TempDir() + "progress"
	mErr := os.MkdirAll(recordDir, 0755)
	if mErr != nil {
		log.Println("mkdir for record dir error,", mErr)
		return mErr
	}
	recordPath := filepath.Join(recordDir, etag)
	pr := progressRecord{}
	//==尝试从旧的进度文件中读取进度
	recordFp, openErr := os.Open(recordPath)
	if openErr == nil {
		progressBytes, readErr := ioutil.ReadAll(recordFp)
		if readErr == nil {
			mErr := json.Unmarshal(progressBytes, &pr)
			if mErr == nil {
				// 检查context 是否过期，避免701错误
				for _, item := range pr.Progresses {
					if storage.IsContextExpired(item) {
						log.Println(item.ExpiredAt)
						pr.Progresses = make([]storage.BlkputRet, storage.BlockCount(fileSize))
						break
					}
				}
			}
		}
		err := recordFp.Close()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if len(pr.Progresses) == 0 {
		pr.Progresses = make([]storage.BlkputRet, storage.BlockCount(fileSize))
	}

	//配置
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	resumeUploader := storage.NewResumeUploader(&cfg)
	ret := storage.PutRet{}
	progressLock := sync.RWMutex{}
	putExtra := storage.RputExtra{
		Progresses: pr.Progresses,
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			progressLock.Lock()
			progressLock.Unlock()
			//将进度序列化，然后写入文件
			pr.Progresses[blkIdx] = *ret
			progressBytes, _ := json.Marshal(pr)
			wErr := ioutil.WriteFile(recordPath, progressBytes, 0644)
			if wErr != nil {
				log.Println("write progress file error,", wErr)
			}
		},
	}
	err := resumeUploader.PutFile(context.Background(), &ret, upKey, etag, localFile, &putExtra)
	if err != nil {
		log.Println(err)
		return err
	}
	//上传成功之后，一定记得删除这个进度文件
	err = os.Remove(recordPath)
	if err != nil {
		log.Println(err)
		return err
	}
	//通知数据库
	ufr := uploadFinishRes{
		Id: id,
	}
	header := map[string]string{
		"Authorization": Authorization,
	}
	_, err = util.HttpPostJson(CONF.ServerAddr+"/file/uploadFinish", ufr, header)
	if err != nil {
		return err
	}
	return nil
}

/*
func uploadFromProgress(recordPath,upKey, etag, localFile string) error {
	progressRecord := ProgressRecord{}
	//==尝试从旧的进度文件中读取进度
	recordFp, openErr := os.Open(recordPath)
	if openErr == nil {
		progressBytes, readErr := ioutil.ReadAll(recordFp)
		if readErr == nil {
			mErr := json.Unmarshal(progressBytes, &progressRecord)
			if mErr == nil {
				// 检查context 是否过期，避免701错误
				for _, item := range progressRecord.Progresses {
					if storage.IsContextExpired(item) {
						log.Println(item.ExpiredAt)
						progressRecord.Progresses = make([]storage.BlkputRet, storage.BlockCount(fileSize))
						break
					}
				}
			}
		}
		err := recordFp.Close()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if len(progressRecord.Progresses) == 0 {
		progressRecord.Progresses = make([]storage.BlkputRet, storage.BlockCount(fileSize))
	}

	//配置
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	resumeUploader := storage.NewResumeUploader(&cfg)
	ret := storage.PutRet{}
	progressLock := sync.RWMutex{}
	putExtra := storage.RputExtra{
		Progresses: progressRecord.Progresses,
		Notify: func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			progressLock.Lock()
			progressLock.Unlock()
			//将进度序列化，然后写入文件
			progressRecord.Progresses[blkIdx] = *ret
			progressBytes, _ := json.Marshal(progressRecord)
			wErr := ioutil.WriteFile(recordPath, progressBytes, 0644)
			if wErr != nil {
				log.Println("write progress file error,", wErr)
			}
		},
	}
	err := resumeUploader.PutFile(context.Background(), &ret, upKey, etag, localFile, &putExtra)
	if err != nil {
		log.Println(err)
		return err
	}
	//上传成功之后，一定记得删除这个进度文件
	err = os.Remove(recordPath)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
*/

//重启上传列表
func UploadRestartTask(Authorization string) (interface{}, error) {
	//进度文件夹
	recordDir := os.TempDir() + "progress"
	//检查文件是否存在
	_, err := os.Stat(recordDir) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return nil, err
		}
	}
	//获取所有文件
	files, _ := ioutil.ReadDir(recordDir)
	if len(files) == 0 {
		return nil, fmt.Errorf("no progress")
	}
	etagArr := make([]string, 0)
	for _, v := range files {
		etagArr = append(etagArr, v.Name())
	}

	//请求未完成列表
	type uploadingParam struct {
		Etags []string `json:"etags"`
	}
	header := map[string]string{
		"Authorization": Authorization,
	}
	data := uploadingParam{
		Etags: etagArr,
	}
	body, err := util.HttpPostJson(CONF.ServerAddr+"/file/uploading", data, header)
	if err != nil {
		return nil, err
	}

	//解析body
	type item struct {
		Id    string `json:"id"`
		Local string `json:"local"`
		Etag  string `json:"etag"`
	}
	type resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data []item
	}
	r := new(resp)
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	if len(r.Data) < 0 {
		return nil, fmt.Errorf("no file")
	}

	//获取上传凭证
	upKey, err := getUpKey(Authorization)
	if err != nil {
		return nil, err
	}

	//TODO 这里有严重高并发问题，大量未上传文件会占用过多线程和网络io和磁盘io，希望七牛内部有解决！！
	for _, s := range r.Data {
		//localFile, upKey, etag, id, Authorization
		go uploadFile(s.Local, upKey, s.Etag, s.Id, Authorization)
	}

	return r.Data, nil
}

//获取上传进度
func LocalUpLoadingProgress() (string, error) {
	recordDir := os.TempDir() + "progress"
	_, err := os.Stat(recordDir)
	if err != nil {
		log.Println(err)
		return "", err
	}
	fileInfoList, err := ioutil.ReadDir(recordDir)
	progressMap := make(map[string]string, 0)
	for _, v := range fileInfoList {
		etag := v.Name()
		//读取文件
		recordPath := filepath.Join(recordDir, etag)
		progressBytes, err := ioutil.ReadFile(recordPath)
		if err != nil {
			log.Println(err)
			continue
		}
		pr := progressRecord{}
		err = json.Unmarshal(progressBytes, &pr)
		if err != nil {
			log.Println(err)
			continue
		}
		total := len(pr.Progresses)
		cum := 0
		for _, item := range pr.Progresses {
			if item.Host != "" {
				cum++
			}
		}

		progressMap[etag] = fmt.Sprintf("%.2f", float32(cum)/float32(total)*100)
	}
	if len(progressMap) > 0 {
		b, err := json.Marshal(progressMap)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return string(b), nil
	}
	return "", nil
}

//获取正在上传的文件
func LocalUpLoadingList() ([]string, error) {
	recordDir := os.TempDir() + "progress"
	_, err := os.Stat(recordDir)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fileInfoList, err := ioutil.ReadDir(recordDir)
	pathArr := make([]string, 0)
	for _, v := range fileInfoList {
		pathArr = append(pathArr, v.Name())
	}
	return pathArr, nil
}
