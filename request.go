package main

import (
	"net/http"
	"io/ioutil"
	"bytes"

	"github.com/yqh231/Urezin/log"
	"encoding/json"
)


type RequestCompose struct {
	request *http.Request
}

func New(method, url string, values interface{}) *RequestCompose{
	params , err := json.Marshal(values)
	if err != nil{
		log.Error.Println(err.Error())
		return nil
	}

	req, er := 	http.NewRequest(method, url, bytes.NewBuffer(params))
	if er != nil{
		log.Error.Println(er.Error())
		return nil
	}

	return &RequestCompose{
		req,
	}
}

func(r *RequestCompose) SetHeader(headers map[string]string){
	for k, v := range headers{
		r.request.Header.Add(k, v)
	}
}


func(r *RequestCompose) ResHandel(res *http.Response, data interface{}){
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		log.Error.Println(err.Error())
		return
	}
	err = json.Unmarshal(body, data)
	if err != nil{
		log.Error.Printf(err.Error())
	}
}


func(r *RequestCompose) GetReq() *http.Request{
	return r.request
}