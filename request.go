package main

import (
	"bytes"
	"encoding/json"
	"github.com/yqh231/Urezin/log"
	"io/ioutil"
	"net/http"
)


type RequestCompose struct {
	request *http.Request
}

func NewReqCompose(method, url string, values interface{}) *RequestCompose{
	var params []byte
	var err error
	if method != "GET" && method != "DELETE"{
		if values != nil{
			params , err = json.Marshal(values)
			if err != nil{
				log.Error.Println(err.Error())
				return nil
			}
		}else{
			params = nil
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

	req, er := 	http.NewRequest(method, url, nil)
	if er != nil{
		log.Error.Println(er.Error())
		return nil
	}

	q := req.URL.Query()
	for k, v := range values.(map[string]string){
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

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