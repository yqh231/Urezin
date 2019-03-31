package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/yqh231/Urezin/log"
)

type RequestCompose struct {
	request *http.Request
}

func NewReqCompose(method, url string, values interface{}) *RequestCompose {
	var params []byte
	var err error
	if method != "GET" && method != "DELETE" {
		if values != nil {
			params, err = json.Marshal(values)
			if err != nil {
				log.Error.Println(err.Error())
				return nil
			}
		}

		req, er := http.NewRequest(method, url, bytes.NewBuffer(params))
		if er != nil {
			log.Error.Println(er.Error())
			return nil
		}
		return &RequestCompose{
			req,
		}

	}

	req, er := http.NewRequest(method, url, nil)
	if er != nil {
		log.Error.Println(er.Error())
		return nil
	}
	if values == nil {
		return &RequestCompose{
			req,
		}
	}
	q := req.URL.Query()
	for k, v := range values.(map[string]string) {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return &RequestCompose{
		req,
	}
}

func (r *RequestCompose) SetHeader(headers map[string]string) {
	for k, v := range headers {
		r.request.Header.Add(k, v)
	}
}

func (r *RequestCompose) ResHandle(res *http.Response, data interface{}) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Error.Println(err.Error())
		switch t := data.(type) {
		case *string:
			*t = string(body)
		case *[]byte:
			*t = []byte(body)
		}
	}
}

func (r *RequestCompose) GetResponseDirect(res *http.Response) string {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error.Println(err.Error())
		return ""
	}
	return string(body)
}

func (r *RequestCompose) GetReq() *http.Request {
	return r.request
}

func Upload(client *http.Client, url string, values map[string]io.Reader) (res *http.Response, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if key == "image" {
			if fw, err = w.CreateFormFile(key, "image"); err != nil {
				return
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		// Add other fields
		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err = client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}
