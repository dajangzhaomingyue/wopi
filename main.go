package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/paulrosania/go-charset/data"
	"io/ioutil"
	whtpp "local/wopi/http"
	wlocal "local/wopi/local"
	"log"
	"net/http"
	"strings"
)

// 定义access_token格式
// 1 type (local 本地处理/http URL处理)
// 2 path (http -> url)/(local -> path)/(oss -> url)
// 3 file_name 文件名（定义显示名）
// 4 version 版本

type Obj struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	Version  string `json:"version"`
}

type fileInfo struct {
	BaseFileName   string `json:"BaseFileName"`
	OwnerId        string `json:"OwnerId"`
	Size           int64  `json:"Size"`
	SHA256         string `json:"SHA256"`
	Version        string `json:"Version"`
	SupportsUpdate bool   `json:"SupportsUpdate,omitempty"`
	UserCanWrite   bool   `json:"UserCanWrite,omitempty"`
	SupportsLocks  bool   `json:"SupportsLocks,omitempty"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/wopi/files/{sub_file_name}", GetFileInfo).Methods(http.MethodGet)
	r.HandleFunc("/wopi/files/{sub_file_name}/contents", GetFileContent).Methods(http.MethodGet)
	r.HandleFunc("/wopi/files/{sub_file_name}/contents", PostFileContent).Methods(http.MethodPost)
	//开启8080端口
	err := http.ListenAndServe(":8080", r)
	log.Println(r)
	if err != nil {
		log.Println("http listen err: ", err)
	}
}

func GetFileInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFileInfo")
	var o Obj
	var err error
	o, err = AccessTokenToObj(r)
	if err != nil {
		log.Println("access to obj err: ", err)
	}
	var info fileInfo
	info.BaseFileName = o.FileName
	info.OwnerId = "admin"
	info.UserCanWrite = true
	info.SupportsLocks = true
	var data []byte
	if o.Type == "local" {
		data, err = wlocal.GetFileData(o.Path)
		if err != nil {
			log.Println("get local data fail: ", err)
			return
		}
	} else if o.Type == "http" {
		data, err = whtpp.GetHttpTest(o.Path)
		if err != nil {
			log.Println("get http data fail: ", err)
			return
		}
	}
	info.Size = int64(len(data))
	info.SHA256, _ = SHA256Byte(data)
	log.Println("debug: sha256_b42: ", info.SHA256)
	info.Version = o.Version
	fmt.Println(info)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(info)
	if err != nil {
		log.Println("encoder file err:", err)
		return
	}
	log.Println("GetFileInfo done...")
}

//获取sha256
func SHA256Byte(buf []byte) (string, error) {
	h := sha256.Sum256(buf)
	return base64.StdEncoding.EncodeToString(h[:]), nil
}

func GetFileContent(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFileContent start.......")
	var o Obj
	var err error
	o, err = AccessTokenToObj(r)
	if err != nil {
		log.Println("access to obj err: ", err)
	}
	var data []byte
	if o.Type == "local" {
		data, err = wlocal.GetFileData(o.Path)
		if err != nil {
			log.Println("get local data fail: ", err)
			return
		}
	} else if o.Type == "http" {
		data, err = whtpp.GetHttpTest(o.Path)
		if err != nil {
			log.Println("get http data fail: ", err)
			return
		}
	}
	w.Header().Set("Content-type", "application/octet-stream")
	_, err = w.Write(data)
	if err != nil {
		log.Println("write file err: ", err)
		return
	}
	log.Println("GetFileContent done !")
}

func PostFileContent(w http.ResponseWriter, r *http.Request) {
	log.Println("PostFileContent start..........")
	var o Obj
	var err error
	o, err = AccessTokenToObj(r)
	if err != nil {
		log.Println("access to obj err: ", err)
	}
	var buf []byte
	buf, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read data err:", err)
		return
	}
	if o.Type == "local" {
		err = wlocal.PostFileData(o.Path, buf)
		if err != nil {
			log.Println("write data err:", err)
			return
		}
	} else if o.Type == "http" {
		//pm := map[string]string{
		//	"path": o.Path,
		//}
		//err = whtpp.PostFormFile(pm, o.PostUrl, "file", o.FileName, bytes.NewReader(buf))
		//if err != nil {
		//	log.Println("post form file err: ", err)
		//	return
		//}
	}

	w.Header().Set("Content-type", "application/octet-stream")
}

func AccessTokenToObj(r *http.Request) (o Obj, err error) {
	val := r.URL.Query()
	tmp, ok := val["access_token"]
	if !ok || len(tmp[0]) == 0 {
		log.Println("access_token not found!")
		return o, errors.New("access_token not found!")
	}
	if len(tmp) == 4 {
		o.Type = tmp[0]
		o.Path = tmp[1]
		o.FileName = tmp[2]
		o.Version = tmp[3]
		//o.PostUrl = tmp[4]
	} else if len(tmp) == 1 {
		oa := strings.Split(tmp[0], ",")
		if len(oa) == 4 {
			o.Type = oa[0]
			o.Path = oa[1]
			o.FileName = oa[2]
			o.Version = oa[3]
			//o.PostUrl = oa[4]
		}
	}
	return o, err
}
