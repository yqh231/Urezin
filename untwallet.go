package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/yqh231/Urezin/log"
)

type UntWallet struct {
	token string
}

func (u *UntWallet) GetToken() {
	res, err := http.Get("http://untwallet.com/account/sms_codes/new")
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Error.Printf("status error %v", res.Status)
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error.Println(err.Error())
	}

	doc.Find("head meta[name=csrf-token]").Each(func(i int, s *goquery.Selection) {
		token, _ := s.Attr("content")
		u.token = token
	})

}

func (u *UntWallet) Work() {

}
