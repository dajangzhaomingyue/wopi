package wlocal

import (
	"io/ioutil"
	"log"
	"os"
)

func GetFileData(fPath string) (data []byte, err error) {
	if fPath != "" {
		var f *os.File
		f, err = os.Open(fPath)
		if err != nil {
			log.Println("open file fail: ", err)
			return data, err
		}
		data, err = ioutil.ReadAll(f)
		if err != nil {
			log.Println("read file fail: ", err)
			return data, err
		}
	}
	return data, err
}

func PostFileData(fPath string, buf []byte) (err error) {
	_ = os.Remove(fPath)
	var f *os.File
	f, err = os.Create(fPath)
	if err != nil {
		log.Printf("open %s err: %s", fPath, err.Error())
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		log.Println("write data err: ", err)
		return
	}
	return nil
}
