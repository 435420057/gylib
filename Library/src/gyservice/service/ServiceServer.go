package service

import (
	"gyservice/proto"
	"golang.org/x/net/context"
	"gylogger"
	"encoding/json"
	"gyservice/respcode"
	cMsg "gycache/message"
	"reflect"
)

func NewServiceServer() (server *ServiceServer) {
	server = &ServiceServer{}
	server.Router = make(map[int32]func(map[string]interface{}) *cMsg.Response)
	return
}

type ServiceServer struct {
	Router map[int32]func(map[string]interface{}) *cMsg.Response
}

func (server *ServiceServer) Serve(ctx context.Context, req *proto.Request) (resp *proto.Response, err error) {

	defer func() {
		if err := recover(); err != nil {
			logger.Infof("server error: %v", err)
		}
	}()

	var cacheContent []byte
	msg, err := cMsg.GetMsg(req.Key, reflect.TypeOf(cMsg.Request{}))
	request := msg.(cMsg.Request)
	if err != nil {
		return
	}
	logger.Debug("Received request:", request)
	call := server.Router[request.Action]
	if call == nil {
		cacheContent, _ = json.Marshal(respcode.RC_GENERAL_SYS_ERR)
	} else {
		logger.Debug("Calling function %s.", request.Action)
		response := call(request.Params)
		cacheContent, _ = json.Marshal(response)
		logger.Debugf("Response from function %s.", string(cacheContent))
	}
	respKey, err := cMsg.CacheMsg(cacheContent)
	resp = &proto.Response{Key:respKey}
	return
}

func (server *ServiceServer) RegHandler(action int32, function func(map[string]interface{}) *cMsg.Response) {
	server.Router[action] = function
}