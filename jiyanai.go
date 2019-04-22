package main

import (
	"fmt"
	"time"

	"github.com/yqh231/Urezin/log"
)

type JiYanAi struct {
	userName string
	passWord string
	url      string
	fileName string
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
		"",
		"60",
	}
}

func (jy *JiYanAi) GetPicture(url string, headers Header) error{
	req := Requests()
	resp, err := req.Get(url, headers)
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	fileName := fmt.Sprintf("/tmp/tm_%s.jpg", time.Now().String()) 
	jy.fileName = fileName
	return resp.SaveFile(fileName)
}

func (jy JiYanAi) Distinguish(file *[]byte, typeId string) string {
	result := new(DistinguishRes)
	req := Requests()
	params := Datas{
		"username": jy.userName,
		"password": jy.passWord,
		"timeout":  jy.timeout,
		"typeid":   typeId,
	}
	resp, er := req.Post(jy.url, params, Files{"image": jy.fileName}) 
	if er != nil {
		log.Error.Println(er.Error())
	}

	err := resp.Json(result)
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
