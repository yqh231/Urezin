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

	reqCompose := NewReqCompose("GET", ym.url,
		map[string]string{
			"action":   "login",
			"username": name,
			"password": password,
		})
	resp, err := ym.client.Do(reqCompose.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	reqCompose.ResHandle(resp, &res)
	stringSlice := strings.Split(res, "|")
	if stringSlice[0] != "success" {
		log.Error.Printf("yi ma get token failed, details is %v", stringSlice[1])
		return
	}
	ym.token = stringSlice[1]

}

func (ym *YiMaCli) GetPhone(itemId, phone string) string {
	var res string

	params := map[string]string{
		"action": "getmobile",
		"token":  ym.token,
		"itemid": itemId,
	}
	if phone != "" {
		params["mobile"] = phone
	}
	reqCompose := NewReqCompose("GET", ym.url,
		params)
	resp, err := ym.client.Do(reqCompose.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	reqCompose.ResHandle(resp, &res)
	stringSlice := strings.Split(res, "|")
	if stringSlice[0] != "success" {
		log.Error.Printf("yi ma get mobile failed, details is %v", stringSlice[1])
		return ""
	}
	return stringSlice[1]

}

func (ym *YiMaCli) GetMessage(mobile, itemId, ifRelease string, callNum int) string {
	var res string

	reqCompose := NewReqCompose("GET", ym.url, map[string]string{
		"action":  "getsms",
		"token":   ym.token,
		"itemid":  itemId,
		"mobile":  mobile,
		"release": ifRelease,
	})
	resp, err := ym.client.Do(reqCompose.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	reqCompose.ResHandle(resp, &res)
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
