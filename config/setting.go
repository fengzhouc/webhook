package config

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	Config *GlobalSetting
	once   = &sync.Once{} //保障线程安全
)

type GlobalSetting struct {
	ServerSetting        ServerSetting        `yaml:"server"`
	WebHookServerSetting WebHookServerSetting `yaml:"webhook"`
	DbSetting            DbSetting            `yaml:"db"`
	CronSetting          CronSetting          `yaml:"cron"`
}

// webhook服务器的基本信息
type ServerSetting struct {
	BaseUrl string `yaml:"baseurl"` // 对外访问的url
	Addr    string `yaml:"addr"`    // 运行host:port
}
type CronSetting struct {
	ListenCron string `yaml:"listen"`
}
type DbSetting struct {
	Sqitepath string `yaml:"path"`
}

// 这里是所有支持的webhook接口，也就是适配的各种
type WebHookServerSetting struct {
	Wx string `yaml:"wx"`
}

// 'import config'的时候就会调用,所以用来做初始化,所以可以不用调用GetInstance去获取config对象
func init() {
	Config = getInstance()
}

// 获取globalSetting对象，单例模式
func getInstance() *GlobalSetting {
	once.Do(func() {
		Config = &GlobalSetting{}
		loadYml() //加载本地配置文件

	})
	return Config
}

// 加载yml中的配置
func loadYml() error {
	// 1. 读取配置文件内容，将返回一个[]byte的内容
	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		// log.Println("[loadYml] ", err)
		return err
	}

	// 2. 使用yaml包进行反序列化
	err = yaml.Unmarshal(file, Config)
	if err != nil {
		log.Print("Unmarshal: ", err)
		return err
	}
	return nil
}
