package http

import (
	"net/http"
	"io/ioutil"
	"redis2server/log"
	"encoding/json"
	"net/url"
)


type Http struct {
	conf  *HttpConf
	value string
}

func (this *Http) SetValue(value string) {
	this.value = value
}

func (this *Http) Http() int32 {
	var code int32
	for i := 0; i < 3; i++ {
		mapResp := this.Post()
		if mapResp!=nil {
			code =  mapResp["code"]
			if code == 200 {
				log.Info("http received success. the code is:", mapResp["code"])
				return code
			}
			log.Info("http received code error.code:", mapResp["code"])
		}

	}
	return code

}

func (this *Http)Post() map[string]int32 {
	//post请求
	data := make(url.Values)
	data["key"] = []string{this.value}
	resp, err := http.PostForm(this.conf.Url,data)
	if err != nil {
		log.Error("http post to phpService failed:", err)
		return nil
	}
	log.Info("http post to phpService success")

	//读取返回值
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("http read resp failed:", err)
	}

	var body map[string]int32
	err =json.Unmarshal(bodyByte, &body)
	if err != nil {
		log.Error("unmarshal json error:",err)
	}
	log.Info("http read resp success:",body)
	return body
}
