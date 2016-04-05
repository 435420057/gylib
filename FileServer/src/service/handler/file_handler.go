package handler

import (
	"gycache/message"
	"gyparam/objId"
	"gyservice/respcode"
	"db/dao"
	"gycache/file"
	"gylogger"
)

func DeleteFileHandler(req map[string]interface{}) (resp *message.Response) {
	resp = message.NewResponse()

	fileId, okFileId := objId.GetObjectIdHexStringWithKey(req, "fileId")

	if !okFileId {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		resp.SetParam("fileId", okFileId)
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

	fileId, okFileId := objId.GetObjectIdHexStringWithKey(req, "fileId")

	if !okFileId {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		resp.SetParam("fileId", okFileId)
		return
	}

	if !file.RefreshExpire(fileId, 1) {
		name, contentType, content, err := dao.LoadFile(fileId)
		if err == nil {
			file.NewCacheFileWithKey(fileId, name, contentType, content, 1)
		} else {
			resp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
			logger.Info(err)
			return
		}
	}

	resp.SetRespCode(respcode.RC_GENERAL_SUCC)
	resp.SetParam("fileId", fileId)
	logger.Debug("LoadFileHandler resp:", resp)

	return
}

func SaveFileHandler(req map[string]interface{}) (resp *message.Response) {
	resp = message.NewResponse()

	fileId, okFileId := objId.GetObjectIdHexStringWithKey(req, "fileId")

	if !okFileId {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		resp.SetParam("fileId", okFileId)
		return
	}

	name, contentType, content, exists := file.GetCacheFile(fileId, true)

	if exists {
		collection, fileId, dbName, err := dao.SaveFile(name, contentType, content)
		logger.Debug("save cached file into mongodb", fileId, err)
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
		logger.Info("Could not find cached file, or expired.")
	}

	return
}

