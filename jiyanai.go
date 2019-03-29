package main

import (
	"github.com/yqh231/Urezin/log"
	"net/http"
)

type JiYanAi struct {
	passWord string
	userName string
	url string
	client *http.Client
	timeout int64
}

type DistinguishRes struct {
	Code int64
	Data struct{
		Id string
		Result string
	}
}

func NewJiYan() *JiYanAi{
	return &JiYanAi{
		"yqh231",
		"2316678",
		"http://aiapi.c2567.com/api/create",
		&http.Client{},
		90,
	}
}

func(jy JiYanAi) GetPicture(url string) *[]byte{
	var res []byte
	req := NewReqCompose("GET", url, nil)
	resp, err := jy.client.Do(req.GetReq())
	if err != nil{
		log.Error.Println(err.Error())
		return nil
	}
	req.ResHandel(resp, &res)
	return &res
}

func(jy JiYanAi) Distinguish(file *[]byte, typeId string) string{
	result := new(DistinguishRes)
	req := NewReqCompose("POST", jy.url,
						map[string][]byte{
							"username": []byte(jy.userName),
							"password": []byte(jy.passWord),
							"url": []byte(jy.url),
							"typeid": []byte(typeId),
							"image": *file,
						})
	resp, err := jy.client.Do(req.GetReq())
	if err != nil{
		log.Error.Println(err.Error())
		return ""
	}
	req.ResHandel(resp, result)
	if result.Code != 10000{
		log.Error.Println(result.Data.Result)
		return ""
	}
	return result.Data.Result
}
