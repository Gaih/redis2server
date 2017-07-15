package http

import (
	"os"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"redis2server/log"
)

type HttpConf struct {
	Url string
}


func NewHttp(path string,name string) *Http{
	mHttp := new(Http)
	conf:=openConf(path,name)
	mHttp.conf = conf
	log.Info("Get redisConfig success")
	return mHttp
}

func openConf(path string, name string) *HttpConf{
	conf := new(HttpConf)
	//读取redis配置文件
	file, err := os.Open(path + name + ".conf")
	if err != nil {
		log.Error("open httpConf error ", err)
		return nil
	}
	fcontent, err := ioutil.ReadAll(file);
	if err != nil {
		log.Error("ReadHttpConf error ", err)
		return nil
	}
	//解析配置文件
	if err = toml.Unmarshal(fcontent, conf); err != nil {
		log.Error("toml.Unmarshal error ", err)
		return nil
	}
	log.Info("read redisConfig success")
	return conf
}