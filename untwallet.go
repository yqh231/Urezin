package main

import (
	"net/http"
	"time"
	"fmt"
	"strconv"
	"io/ioutil"
	"strings"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/yqh231/Urezin/log"
)

type UntWallet struct {
	token string
	cookie string

	jiyan *JiYanAi
	yima *YiMaCli
}

func NewUntWallet() *UntWallet{
	return &UntWallet{
		jiyan: NewJiYan(),
		yima: NewYiMa(),
	}
}

func (u *UntWallet) GetToken() {
	req := NewReqCompose("GET", "http://untwallet.com/account/sms_codes/new", nil)
	res, err := u.jiyan.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	u.GetCookie(res)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error.Printf("status error %v", res.Status)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error.Println(err.Error())
		return
	}

	doc.Find("head meta[name=csrf-token]").Each(func(i int, s *goquery.Selection) {
		token, _ := s.Attr("content")
		u.token = token
	})

}

func (u *UntWallet) GetCookie(resp *http.Response){
	cookie, ok := resp.Header["Set-Cookie"]
	if !ok{
		return
	}
	c_slice := strings.Split(cookie[0], ";")
	u.cookie = strings.TrimSpace(c_slice[0])
}


func (u *UntWallet) GetPicture(url string) *[]byte{
	var res []byte
	req := NewReqCompose("GET", url, nil)
	req.SetHeader(map[string]string{
		"Cookie": u.cookie,
	})

	resp, err := u.jiyan.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	u.GetCookie(resp)
	req.ResHandle(resp, &res)
	return &res
}

func (u *UntWallet) VerifyCode() string{
	var result string

PHONE_DUPLICATE:
	u.yima.Login("yqh231", "2316678")
	phone := u.yima.GetPhone("21714", "")
	u.GetToken()

CODE_ERR:
	now := time.Now().UnixNano() / int64(time.Millisecond)
	image := u.GetPicture("http://untwallet.com/rucaptcha/?" + strconv.Itoa(int(now)))
	captcha := u.jiyan.Distinguish(image, "3000")
	ioutil.WriteFile("t.jpg", *image, 0644)
	v := url.Values{}
	v.Add("utf8", "&#x2713")
	v.Add("user[area_code]", "86")
	v.Add("user[number]", phone)
	v.Add("_rucaptcha", captcha)
	v.Add("commit", "发送验证码")
	req := NewReqCompose("POST", "http://untwallet.com/account/sms_codes", 
		v.Encode())
	req.SetHeader(map[string]string{"Cookie": u.cookie, "X-CSRF-Token": u.token})
	resp, err := u.jiyan.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	u.GetCookie(resp)
	defer resp.Body.Close()

	req.ResHandle(resp, &result)
	if result != ""{
		log.Error.Println(result)
		if strings.HasPrefix(result, "alert('验证码不正确')"){
			goto CODE_ERR
		}
		if strings.HasPrefix(result, "alert"){
			goto PHONE_DUPLICATE 
		}
	}
	fmt.Println("return")
	return u.yima.GetMessage(phone, "21714", "0", 0)
}


func (u *UntWallet) Validate(code string){
	var result string
	code_int, _ := strconv.Atoi(code)
	req := NewReqCompose("POST", "http://untwallet.com/account/sms_codes/validate", 
		map[string]interface{}{
			"utf8": "&#x2713;",
			"_method": "patch",
			"sms_validation_code[value]": code_int,
			"commit": "验证"})
	req.SetHeader(map[string]string{"Cookie": u.cookie})
	resp, err := u.jiyan.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	u.GetCookie(resp)
	defer resp.Body.Close()
	req.ResHandle(resp, &result)
	if result != ""{
		log.Error.Println(result)
		return
	}
}

func (u *UntWallet) RegisterAccount(){
	var result string

	req := NewReqCompose("POST", "http://untwallet.com/account",
		map[string]interface{}{
			"utf8": "&#x2713;",
			"user[password]": 2316678,
			"user[password_confirmation]": 2316678,
			"user[area_code]": 86,
			"user[number]": 18689468297,
			"commit": "注册",})
	req.SetHeader(map[string]string{"Cookie": u.cookie})
	resp, err := u.jiyan.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	u.GetCookie(resp)
	defer resp.Body.Close()
	req.ResHandle(resp, &result)
	if result != ""{
		log.Error.Println(result)
		return
	}
}

func (u *UntWallet) Transaction(){
	var result string
	req := NewReqCompose("POST", "http://untwallet.com/my/transactions",
		map[string]interface{}{
			"utf8": "&#x2713;",
			"transaction[lockde]": false,
			"transaction[receiver_area_code]": 86,
			"transaction[receiver_number]": 18689468297,
			"transaction[amount]": 3,
			"transaction[password]": 2316678,
			"commit": "转出",
		})
	req.SetHeader(map[string]string{"Cookie": u.cookie})
	resp, err := u.jiyan.client.Do(req.GetReq())
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	u.GetCookie(resp)
	defer resp.Body.Close()
	req.ResHandle(resp, &result)
	if result != ""{
		log.Error.Println(result)
		return
	}
}


func (u *UntWallet) Work() {
	code := u.VerifyCode()
	u.Validate(code)
	u.RegisterAccount()
	u.Transaction()
}
