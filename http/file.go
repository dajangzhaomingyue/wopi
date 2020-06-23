package whtpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

func GetHttpTest(baseUrl string) (buf []byte, err error) {
	client := &http.Client{}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, baseUrl, nil); err != nil {
		log.Println("new request fail: ", err)
		return buf, err
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		log.Println("client do fail: ", err)
		return buf, err
	}
	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("read request body fail: ", err)
		return buf, err
	}
	return buf, nil
}

func PostFormFile(pM map[string]string, baseUrl string, fieldName, fileName string, fileReader io.Reader) error {
	var err error
	var b []byte
	if b, err = json.Marshal(pM); err != nil {
		log.Println("marshal pm err: ", err)
		return err
	}

	bodyBuf := bytes.NewBuffer(b)

	bodyWriter := multipart.NewWriter(bodyBuf)

	//form file
	var formFile io.Writer
	if formFile, err = bodyWriter.CreateFormFile(fieldName, fileName); err != nil {
		log.Println("create form file err: ", err)
		return err
	}

	if _, err = io.Copy(formFile, fileReader); err != nil {
		log.Println("io copy fileReader err: ", err)
		return err
	}

	// 发送表单
	contentType := bodyWriter.FormDataContentType() // 表单类型
	err = bodyWriter.Close()                        // 发送之前必须调用Close()以写入结尾行
	if err != nil {
		log.Println("body writer close err: ", err)
		return err
	}

	var rep *http.Response
	if rep, err = http.Post(baseUrl, contentType, bodyBuf); err != nil {
		log.Println("post http request err: ", err)
		return err
	}
	if rep.StatusCode != 200 {
		return err
	}
	fmt.Println(rep.StatusCode)
	return nil
}
