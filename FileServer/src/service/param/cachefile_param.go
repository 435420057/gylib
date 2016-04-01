package param

import (
	"github.com/sluu99/uuid"
	log "github.com/kyugao/go-logger/logger"
)

func GetCacheKey(req map[string]interface{}) (fileId string, ok bool) {
	tempFileId, ok := req["fileId"]
	if ok {
		fileId = tempFileId.(string)
	} else {
		return
	}
	_, err := uuid.FromStr(fileId)
	if err != nil {
		log.Debugf("conver cache file key string %s, err %v", fileId, err)
		ok = false
	} else {
		log.Debugf("conver cache file key string %s done.", fileId)
		ok = true
	}
	return
}

func GetFileId(req map[string]interface{}) (fileId string, ok bool) {
	tempFileId, ok := req["fileId"]
	if ok {
		fileId = tempFileId.(string)
	}
	return
}