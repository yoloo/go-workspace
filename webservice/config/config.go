package config

import (
	"encoding/xml"
	"io/ioutil"
	"github.com/go-redis/redis"
	"webservice/util"
)

var GlobalConfig = &AppConfig{}

var MyCache *redis.Client = nil
var MyRedis *redis.Client = nil

func LoadConfig() error {
	content,err := ioutil.ReadFile(util.GetCurrentPath() + "appconfig.xml")
	if err != nil {
		return err
	}

	err = xml.Unmarshal(content, GlobalConfig)
	if err != nil {
		return err
	}

	MyCache = redis.NewClient(&redis.Options{
		Addr:     GlobalConfig.Mycache.Address,
		Password: GlobalConfig.Mycache.Password, // no password set
		DB:       GlobalConfig.Mycache.Db,        // use default DB
	})

	MyRedis = redis.NewClient(&redis.Options{
		Addr:     GlobalConfig.Myredis.Address,
		Password: GlobalConfig.Myredis.Password, // no password set
		DB:       GlobalConfig.Myredis.Db,        // use default DB
	})
	return nil
}



