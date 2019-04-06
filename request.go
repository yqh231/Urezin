package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/yqh231/Urezin/log"
)

type RequestCompose struct {
	request *http.Request
}

func NewReqCompose(method, url string, values interface{}) *RequestCompose {

	if method != "GET" && method != "DELETE" {
		req, _ := http.NewRequest(method, url, strings.NewReader(values.(string)))
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
	r.request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Mobile Safari/537.36")
	r.request.Header.Add("Refer", "http://untwallet.com/account/sign_in")
	r.request.Header.Add("Host", "untwallet.com")
	r.request.Header.Add("Accept", "text/javascript, application/javascript, application/ecmascript, application/x-ecmascript, */*; q=0.01")
}

func (r *RequestCompose) ResHandle(res *http.Response, data interface{}) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	err = json.Unmarshal(body, data)
	if err != nil {
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
