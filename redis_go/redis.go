package redis_go

import (
	"github.com/garyburd/redigo/redis"
	"redis2server/log"
	"encoding/json"
	"time"
)

type Redis struct {
	conf *RedisConf
	conn redis.Conn
	pool *redis.Pool
}
//
//func (this *Redis) NewRedisConn() {
//
//	//log.Info("typeIP:", this.conf.Type, this.conf.IP)
//	dialPass := redis.DialPassword(this.conf.Password)
//	conn, err := redis.Dial(this.conf.Type, this.conf.IP, dialPass)
//	if err != nil {
//		log.Error("redis connect faild:", err)
//	}
//	//fmt.Println("redis conn success")
//	this.conn = conn
//}

//新建连接池
func (this *Redis) NewPool() {
	this.pool =  &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", this.conf.IP)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", this.conf.Password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}


//向redis添加key和value
func (this *Redis) SetToRedis(key string, value interface{}) {

	_, err := this.conn.Do("SET", key, value)
	if err != nil {
		log.Error("redis set key", key, "failed", err)
		return
	}
	log.Info("redis set ", key, " success")
}

//通过key获取读取数据库
func (this *Redis) GetToRedis(key string) string {

	// 取json数据
	var data map[string]string
	// json数据在go中是[]byte类型，所以此处用redis.Bytes转换
	value, err := redis.Bytes(this.conn.Do("GET", key))
	if err != nil {
		log.Error("Get redis data error:",key, err)
	}
	// 将json转行成string
	if value!=nil {
		errShal := json.Unmarshal(value, &data)
		if errShal != nil {
			log.Error("trans to json err:", key,errShal)
		}
		//fmt.Println( data["hello"])
		return data["hello"]
	}
	return ""
}

//获取指定条件的keys
func (this *Redis) GetKeys(title string) []string {
	value, err := redis.Strings(this.conn.Do("keys", title+"*"))
	if err != nil {
		log.Error("Get redis data error:", err)
	}
	return value
}

//删除key和对应数据
func (this *Redis) DelToRedis(key string) {
	_, err := this.conn.Do("DEL", key)
	if err != nil {
		log.Error("redis del key", key, "failed")
		return
	}
	log.Info("redis del ", key, " success")
}

func (this *Redis) Close() {
	this.conn.Close()
}
func (this *Redis) GetConn(){
	this.conn = this.pool.Get()
}
//redis事务处理
func (this *Redis) Multi() {
	this.conn.Do("MULTI")
}
func (this *Redis)Exec(){
	this.conn.Do("EXEC")
}

