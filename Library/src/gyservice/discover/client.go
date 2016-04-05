package discover

import (
	"gylogger"
	"gyservice/proto"
	"gyservice/action"
)

type ServiceRange struct {
	Start int32
	End   int32
	Node  string
}

var pool map[int32]*ServiceRange

func init() {
	pool = make(map[int32]*ServiceRange)
}

func RegisterNode(start int32, serviceRange *ServiceRange) {
	pool[start] = serviceRange
}

func GetClient(actionCode action.Action) (client proto.ServiceClient, serviceNodeName string) {
	code := int32(actionCode)
	for key, val := range pool {
		logger.Debugf("check action code:name %d:%s, get pool element key:%d, start:end:name %d:%d:%s.", code, actionCode, key, val.Start, val.End, val.Node)
		if code >= key && code <= val.End {
			client = GetServiceClient(val.Node)
			serviceNodeName = val.Node
			logger.Debugf("got client %v", client)
			break;
		}
	}
	return
}