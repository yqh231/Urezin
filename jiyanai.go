package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yqh231/Urezin/log"
)

type JiYanAi struct {
	userName string
	passWord string
	url      string
	client   *http.Client
	timeout  string
}

type DistinguishRes struct {
	Code int64
	Data struct {
		Id     string
		Result string
	}
}

func NewJiYan() *JiYanAi {
	return &JiYanAi{
		"yqh231",
		"2316678",
		"http://aiapi.c2567.com/api/create",
		&http.Client{},
		"60",
	}
}

func (jy JiYanAi) GetPicture(url string) *[]byte {
	var res []byte
	req := NewReqCompose("GET", url, nil)
	resp, err := jy.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	req.ResHandle(resp, &res)
	return &res
}

func (jy JiYanAi) Distinguish(file *[]byte, typeId string) string {
	result := new(DistinguishRes)
	params := map[string]io.Reader{
		"image":    bytes.NewReader(*file),
		"username": strings.NewReader(jy.userName),
		"password": strings.NewReader(jy.passWord),
		"timeout":  strings.NewReader(jy.timeout),
		"typeid":   strings.NewReader(typeId),
	}
	resp, er := Upload(jy.client, jy.url, params)
	if er != nil {
		log.Error.Println(er.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	err = json.Unmarshal(body, result)
	fmt.Println(result)
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}

	if result.Code != 10000 {
		log.Error.Println(result.Data.Result)
		return ""
	}
	return result.Data.Result
}
