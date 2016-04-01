package service

import (
	log "github.com/kyugao/go-logger/logger"
	"external/action"
	"service/handler"
	"external/service"
)

func InitServer() *service.ServiceServer {
	log.Info("Register file handlers.")
	server := service.NewServiceServer()
	server.RegHandler(int32(action.Action_LoadFile), handler.LoadFileHandler)
	server.RegHandler(int32(action.Action_SaveFile), handler.SaveFileHandler)
	server.RegHandler(int32(action.Action_DeleteFile), handler.DeleteFileHandler)
	server.RegHandler(int32(action.Action_CheckAppVersion), handler.CheckAppVersion)
	server.RegHandler(int32(action.Action_CheckFirmwareVersion), handler.CheckFirmwareVersion)
	return server
}
