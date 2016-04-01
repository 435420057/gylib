package handler

import (
	"external/action"
	"golang.org/x/net/context"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"external/discover"
	log "github.com/kyugao/go-logger/logger"
	"external/respcode"
	"external/cache/message"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	actionName := vars["action"]
	var respByte []byte

	actionCode, ok := action.ActionFromName(actionName)
	if !ok {
		log.Debug("Unsupport function", actionCode.Name())
		respByte, _ = json.Marshal(respcode.RC_UNSUPPORT_FUNCTION)
	} else {
		r.ParseForm()
		params := make(map[string]interface{})
		for k, _ := range r.Form {
			params[k] = r.FormValue(k)
		}
		request := &message.Request{Action:actionCode, Params:params}

		cachedReq, err := message.CacheReq(request)
		if err != nil {
			respByte, _ = json.Marshal(respcode.RC_GENERAL_SYS_ERR)
		} else {
			log.Debugf("Cached request key %s request: %v", cachedReq.RequestKey, request)
			client, _ := discover.GetClient(actionCode)

			if client == nil {
				respByte, _ = json.Marshal(respcode.RC_SERVICE_UNAVAILABLE)
			} else {
				clientResp, err := client.Serve(context.Background(), cachedReq)
				log.Debug("response from node:", clientResp, err)
				respByte, err = message.GetCacheResp(clientResp.ResponseKey)
				if err != nil {
					log.Debugf("Get cache resp error: %s.", err.Error())
					respByte, _ = json.Marshal(respcode.RC_GENERAL_SYS_ERR)
				} else {
					log.Debugf("Get cache resp: %s.", string(respByte))
				}
			}
		}
	}

	w.Write(respByte)
}