package upload

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// 接收文件
func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	var mf multipart.File
	var data []byte
	fPath := r.FormValue("path")
	mf, _, err := r.FormFile("file")
	if err != nil {
		DealHttpErr(err, w)
	}
	data, err = ioutil.ReadAll(mf)
	if err != nil {
		DealHttpErr(err, w)
	}
	// 打开文件
	_ = os.Remove(fPath)
	var f *os.File
	f, err = os.Open(fPath)
	if err != nil {
		DealHttpErr(err, w)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		DealHttpErr(err, w)
	}
	_, err = w.Write([]byte("write success"))
	if err != nil{
		DealHttpErr(err, w)
	}
	w.WriteHeader(http.StatusOK)
}

// 获取文件
func LoadFile(w http.ResponseWriter, r *http.Request) {
	var obj struct {
		FilePath string `json:"file_path"`
	}
	var err error

	if err = json.NewDecoder(r.Body).Decode(obj); err != nil {
		DealHttpErr(err, w)
	}
	var f *os.File
	f, err = os.Open(obj.FilePath)
	if err != nil {
		DealHttpErr(err, w)
	}
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		DealHttpErr(err, w)
	}
	_, fileName := filepath.Split(obj.FilePath)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(data)
	if err != nil{
		DealHttpErr(err, w)
	}
}

func DealHttpErr(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return
}
