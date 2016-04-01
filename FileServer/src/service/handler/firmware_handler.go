package handler

import (
	"external/cache/message"
	"service/param"
	"external/respcode"
	"db/dao"
	"github.com/kyugao/go-logger/logger"
)

func CheckFirmwareVersion(req map[string]interface{}) (resp *message.Response) {
	resp = message.NewResponse()

	version, okVersion := param.GetFirmwareVersion(req)
	if !(okVersion) {
		resp.SetRespCode(respcode.RC_GENERAL_PARAM_ERR)
		return
	}
	logger.Info(version)

	isLatest, latestVersion, fileId, err := dao.IsLatestFirmwareVersion(version)
	if err != nil {
		resp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
		resp.SetParam("error", err.Error())
	} else {
		resp.SetRespCode(respcode.RC_GENERAL_SUCC)
		resp.SetParam("latest", isLatest)
		if !isLatest {
			resp.SetParam("latest_version", latestVersion)
			resp.SetParam("fileId", fileId)
		}
	}
	return
}