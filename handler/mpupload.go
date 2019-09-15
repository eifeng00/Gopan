package handler

import (
	"filestore-server/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	rPool "filestore-server/cache/redis"
	dblayer "filestore-server/db"

	"github.com/garyburd/redigo/redis"
)

// MultipartUploadinfo : 初始化信息
type MultipartUploadinfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// InitialMultipartUploadHandler : 初始化分块上传操作
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))

	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invaild", nil).JSONBytes())
		return
	}

	//2. 获取redis的一个链接
	rConn := rPool.RedisPool().Get()

	defer rConn.Close()
	//3. 生成分块上传的初始化信息

	upInfo := MultipartUploadinfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}
	//4. 将初始化信息写入redis缓存
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSTS", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)
	//5. 将初始化的信息返回给客户端
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析用户请求参数
	r.ParseForm()

	_ = r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	//2. 获取redis连接池的一个连接
	rConn := rPool.RedisPool().Get()

	defer rConn.Close()

	//3. 获得文件据并, 用于存储分块内容
	fpath := "/tmp/filestores_cache/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)

	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		fmt.Printf(err.Error())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
	}

	//4. 更新redis缓存状态
	rConn.Do("HSTS", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	//5. 返回处理的结果给客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析请求参数
	r.ParseForm()

	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")
	//2. 获取redis连接池的中的一个连接
	rConn := rPool.RedisPool().Get()

	defer rConn.Close()
	//3. 通过uploadid 查询redis,判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "completer update failed", nil).JSONBytes())
		return
	}

	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}

	//TODO : 4. 合并分块

	//5. 更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	//6. 相应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
