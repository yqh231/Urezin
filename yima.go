package main

import (
	"net/http"
	"strings"
	"time"
	"fmt"

	"github.com/yqh231/Urezin/log"
)

type YiMaCli struct {
	token  string
	client *http.Client
	url    string
}

func NewYiMa() *YiMaCli {
	return &YiMaCli{
		client: &http.Client{},
		url:    "http://api.fxhyd.cn/UserInterface.aspx",
	}
}

func (ym *YiMaCli) Login(name, password string) {
	var res string

	req := Requests()

	resp, err := req.Get(ym.url, Params{"action": "loginIn", "name": name, "password": password})
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	if resp != nil {
		return
	}

	res = resp.Text()
	stringSlice := strings.Split(res, "|")
	if stringSlice[0] != "success" {
		log.Error.Printf("yi ma get token failed, details is %v", stringSlice[1])
		return
	}
	ym.token = stringSlice[1]

}

func (ym *YiMaCli) GetPhone(itemId, phone string) string {
	var res string
	req := Requests()
	params := Params{
		"action": "getmobile",
		"token": ym.token,
		"itemid": itemId,
	}
	if phone != "" {
		params["mobile"] = phone
	}

	resp, err := req.Get(ym.url, params) 
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	res = resp.Text()
	stringSlice := strings.Split(res, "|")
	if stringSlice[0] != "success" {
		log.Error.Printf("yi ma get mobile failed, details is %v", stringSlice[1])
		return ""
	}
	return stringSlice[1]

}

func (ym *YiMaCli) GetMessage(mobile, itemId, ifRelease string, callNum int) string {
	var res string

	req := Requests()
	resp, err := req.Get(ym.url, Params{"action": "getsms", "token": ym.token, "itemid": itemId, "mobile": mobile, "release": ifRelease})
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}

	res = resp.Text()
	stringSlice := strings.Split(res, "|")
	fmt.Println(res)
	if !strings.HasPrefix(res, "success") && callNum < 20 {
		log.Error.Printf("yi ma get message fail, details is %v", stringSlice[1])
		ym.GetMessage(mobile, itemId, ifRelease, callNum+1)
		time.Sleep(5 * time.Second)
	}

	if callNum >= 20 {
		return ""
	}
	return stringSlice[1]
}
