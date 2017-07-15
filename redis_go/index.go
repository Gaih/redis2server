package redis_go

import (
	"redis2server/log"
	"os"
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

type RedisConf struct {
	Type     string
	IP       string
	Username string
	Password string
}

func InitRedis( confPath string, name string) *Redis {
	redis := NewRedis(confPath, name)
	return redis
}

func NewRedis(path string, name string) *Redis {
	redis := new(Redis)
	conf := openConf(path, name)
	redis.conf = conf
	log.Info("Get redisConfig success")
	return redis
}
func openConf(path string, name string) *RedisConf {
	conf := new(RedisConf)
	//读取redis配置文件
	file, err := os.Open(path + name + ".conf")
	if err != nil {
		log.Error("open redisConf error ", err)
		return nil
	}
	fcontent, err := ioutil.ReadAll(file);
	if err != nil {
		log.Error("ReadRedisConf error ", err)
		return nil
	}
	//解析配置文件
	if err = toml.Unmarshal(fcontent, conf); err != nil {
		log.Error("toml.Unmarshal error ", err)
		return nil
	}
	log.Info("read redisConfig success", *conf)
	return conf
}
