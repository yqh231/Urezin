package main

import (
	"github.com/yqh231/Urezin/log"
	"net/http"
	"strings"
	"time"
)


type HuoYunCli struct {
	token string
	client *http.Client
	url string
	getMsgFail chan map[string]string
}


func NewHuoYun() *HuoYunCli{
	return &HuoYunCli{
		client: &http.Client{},
		url: "http://huoyun888.cn/api/do.php",
		getMsgFail: make(chan map[string]string, 20),
	}
}

func(hy *HuoYunCli) Login(name, password string){
	var res string

	reqCompose := NewReqCompose("GET", hy.url,
		map[string]string{
		"action": "loginIn",
		"name": name,
		"password": password,
	})
	resp, err := hy.client.Do(reqCompose.GetReq())
	if err != nil{
		log.Error.Println(err.Error())
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		return
	}

	reqCompose.ResHandel(resp, &res)
	stringSlice := strings.Split("|", res)
	if stringSlice[0] == "0"{
		log.Error.Printf("login huo yun fail, details is %v", stringSlice[1])
	}
	hy.token = stringSlice[1]

}


func(hy *HuoYunCli) GetPhone(sid, phone string) string{
	var res string

	params := map[string]string{
		"action": "getPhone",
		"token": hy.token,
		"sid": sid,
	}

	if phone != ""{
		params["phone"] = phone
	}

	reqCompose := NewReqCompose("GET", hy.url, params)
	resp, err := hy.client.Do(reqCompose.GetReq())
	if err != nil{
		log.Error.Panicln(err.Error())
		return ""
	}
	if resp != nil {
		defer resp.Body.Close()
		return ""
	}

	reqCompose.ResHandel(resp, &res)
	stringSlice := strings.Split("|", res)
	if stringSlice[0] == "0"{
		log.Error.Printf("huo yun get phone fail, details is %v", stringSlice[1])
	}
	return stringSlice[1]
}

func(hy *HuoYunCli) GetMessage(sid, phone, author string) string{
	var res string

	params := map[string]string{
		"action": "getMessage",
		"token": hy.token,
		"sid": sid,
		"phone": phone,
		"author": author,
	}

	if phone != ""{
		params["phone"] = phone
	}

	reqCompose := NewReqCompose("GET", hy.url, params)
	resp, err := hy.client.Do(reqCompose.GetReq())
	if err != nil{
		log.Error.Panicln(err.Error())
		return ""
	}
	if resp != nil {
		defer resp.Body.Close()
		return ""
	}

	reqCompose.ResHandel(resp, &res)
	stringSlice := strings.Split("|", res)
	if stringSlice[0] == "0"{
		log.Error.Printf("huo yun get message fail, details is %v", stringSlice[1])
		hy.GetMessage(sid, phone, author)
		time.Sleep(3 * time.Second)
	}
	return stringSlice[1]
}
