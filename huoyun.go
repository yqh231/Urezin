package main

import (
	"github.com/yqh231/Urezin/log"
	"net/http"
	"strings"
	"time"
)

type HuoYunCli struct {
	token  string
	client *http.Client
	url    string
}

func NewHuoYun() *HuoYunCli {
	return &HuoYunCli{
		client: &http.Client{},
		url:    "http://huoyun888.cn/api/do.php",
	}
}

func (hy *HuoYunCli) Login(name, password string) {
	var res string
	req := Requests()

	resp, err := req.Get(hy.url, Params{"action": "loginIn", "name": name, "password": password})
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	if resp != nil {
		return
	}

	res = resp.Text()
	stringSlice := strings.Split("|", res)
	if stringSlice[0] == "0" {
		log.Error.Printf("login huo yun fail, details is %v", stringSlice[1])
	}
	hy.token = stringSlice[1]

}

func (hy *HuoYunCli) GetPhone(sid, phone string) string {
	var res string
	req := Requests()

	params := Params{
		"action": "getPhone",
		"token": hy.token,
		"sid": sid,
	}

	if phone != "" {
		params["phone"] = phone
	}
	resp, err := req.Get(hy.url, params) 
	if err != nil {
		log.Error.Panicln(err.Error())
		return ""
	}
	if resp != nil {
		return ""
	}

	res = resp.Text()
	stringSlice := strings.Split("|", res)
	if stringSlice[0] == "0" {
		log.Error.Printf("huo yun get phone fail, details is %v", stringSlice[1])
	}
	return stringSlice[1]
}

func (hy *HuoYunCli) GetMessage(sid, phone, author string) string {
	var res string

	req := Requests()
	params := Params{
		"action": "getMessage",
		"token":  hy.token,
		"sid":    sid,
		"phone":  phone,
		"author": author,
	}

	if phone != "" {
		params["phone"] = phone
	}

	resp, err := req.Get(hy.url, params) 
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	if resp != nil {
		return ""
	}

	res = resp.Text()
	stringSlice := strings.Split("|", res)
	if stringSlice[0] == "0" {
		log.Error.Printf("huo yun get message fail, details is %v", stringSlice[1])
		hy.GetMessage(sid, phone, author)
		time.Sleep(3 * time.Second)
	}
	return stringSlice[1]
}
