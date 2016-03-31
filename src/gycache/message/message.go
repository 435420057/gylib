package message

import (
	"gycache"
	"gyuuid"
	"encoding/json"
	"reflect"
)

func CheckExists(msgKey string) (bool, error) {
	return cache.Exist(msgKey)
}

func CacheMsg(msg interface{}) (key string, err error) {
	msgBytes, err := json.Marshal(msg)
	key = uuid.Rand().Hex()
	err = cache.Set(key, msgBytes)
	if err == nil {
		cache.Expire(key, 1)
	}
	return
}

func GetMsg(key string, refType reflect.Type) (msg interface{}, err error) {
	msg = reflect.New(refType)
	msgBytes, err := cache.GetBytes(key)
	if err != nil {
		return
	}
	cache.Del(key)
	err = json.Unmarshal(msgBytes, &msg)
	return
}