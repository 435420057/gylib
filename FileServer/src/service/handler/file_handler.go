package handler
import (
	"external/cache/message"
	"service/param"
	"external/respcode"
	"db/dao"
	"external/cache/file"
	log "github.com/kyugao/go-logger/logger"
)

func DeleteFileHandler(req map[string]interface{}) (resp *message.Response) {
	resp = message.NewResponse()

	var fileId string
	var ok bool

	if fileId, ok = param.GetFileId(req); !ok {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		return
	}

	err := dao.DeleteFile(fileId)
	if err == nil {
		resp.SetRespCode(respcode.RC_GENERAL_SUCC)
	} else {
		resp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
	}

	return
}

func LoadFileHandler(req map[string]interface{}) (resp *message.Response) {
	resp = message.NewResponse()
	var fileId string
	var ok bool

	if fileId, ok = param.GetFileId(req); !ok {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		return
	}

	if !file.RefreshExpire(fileId, 1) {
		name, contentType, content, err := dao.LoadFile(fileId)
		if err == nil {
			file.NewCacheFileWithKey(fileId, name, contentType, content, 1)
		} else {
			resp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
			log.Info(err)
			return
		}
	}

	resp.SetRespCode(respcode.RC_GENERAL_SUCC)
	resp.SetParam("fileId", fileId)
	log.Debug("LoadFileHandler resp:", resp)

	return
}

func SaveFileHandler(req map[string]interface{}) (resp *message.Response) {
	resp = message.NewResponse()
	var fileId string
	var ok bool

	if fileId, ok = param.GetCacheKey(req); !ok {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		return
	}

	name, contentType, content, exists := file.GetCacheFile(fileId, true)

	if exists {
		collection, fileId, dbName, err := dao.SaveFile(name, contentType, content)
		log.Debug("save cached file into mongodb", fileId, err)
		if err == nil {
			resp.SetRespCode(respcode.RC_GENERAL_SUCC)
			resp.SetParam("fileId", fileId)
			resp.SetParam("collection", collection)
			resp.SetParam("dbName", dbName)
		} else {
			resp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
		}
	} else {
		resp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
		log.Info("Could not find cached file, or expired.")
	}

	return
}

