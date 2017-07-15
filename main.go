package main

import (
	"sync"
	"fmt"
	"redis2server/redis_go"
	"time"
	"redis2server/http"
	"redis2server/log"
	"strings"
)

var wait sync.WaitGroup
var mutex sync.Mutex
func main() {
	//log设置保存
	log.NewLogger("")
	//设置定时
	timer()
}

func timer(){
	doRedis()
	timer1:=time.NewTicker(10*time.Second)
	for  {
		select {
		case <-timer1.C:
			doRedis()
		}
	}
}

func doRedis() {
	mRedis := redis_go.InitRedis( "", "redis")
	mRedis.NewPool()
	mRedis.GetConn()
	//获取指定条件的key集合
	value := mRedis.GetKeys("p_111_10008")
	mRedis.Close()
	a := time.Now()
	//遍历所有的keys，除去error的部分
	for _, v := range value {
		if !strings.Contains(v, "error") {
			mRedis.GetConn()
			data := mRedis.GetToRedis(v)
			mRedis.Close()
			fmt.Println("key:",v, ":", data)
			wait.Add(1)
			go goroutine(mRedis,v,data)
		}
	}

	//data:=mRedis.GetToRedis("p_111_1000175")
	//goroutine(mRedis,"p_111_1000175",data)
	wait.Wait()
	fmt.Println(a, "----", time.Now())
}

func goroutine(mRedis *redis_go.Redis,key string,value string) {
	defer wait.Done()
	defer mRedis.Close()
	//读取配置文件
	mHttp:=http.NewHttp("","server")
	//设置从redis读取的值
	mHttp.SetValue(value)
	//发送获取返回值
	code:=mHttp.Http()
	if code ==200{
		log.Info(key,"#success#",time.Now())
	}else {
		//删除之前的key，写入key_error和错误信息
		//互斥锁
		mutex.Lock()
		log.Info("code is error, change redis :",key)
		mRedis.GetConn()
		mRedis.Multi()
		mRedis.DelToRedis(key)
		mRedis.SetToRedis(key+"_error",fmt.Sprint("{\"errorCode:\"",code,"}"))
		mRedis.Exec()
		mutex.Unlock()
	}
}
