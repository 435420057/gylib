package handler
import (
	"net/http"
	log "github.com/kyugao/go-logger/logger"
	"io/ioutil"
	cFile "gycache/file"
	"encoding/json"
	"gycache/token"
	"github.com/gorilla/mux"
	//"external/discover"
	//"external/action"
	"golang.org/x/net/context"
	"fmt"
	"gyservice/respcode"
	"gycache/message"
)

// 5MB
const MAX_MEMORY = 5 * 1024 * 1024

/*
Receive upload files, save all the file information into cache server.
 */
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	agentResp := message.NewResponse()

	if err := r.ParseMultipartForm(MAX_MEMORY); err != nil {
		agentResp.SetRespCode(respcode.RC_GENERAL_SYS_ERR)
		agentResp.SetParam("error", "The file uploaded exceeded the limitation 5M.")
	} else {
		tokenStr := r.MultipartForm.Value["token"][0]

		// validate the tokenStr
		_, _, ok := token.Validate(tokenStr)
		if !ok {
			log.Debugf("token %s expired.\n", tokenStr)
			agentResp.SetRespCode(respcode.RC_GENERAL_APP_ERR)
			agentResp.SetParam("error", "Token invalid or expired.")
		} else {
			// cache file information
			fileHeaders := r.MultipartForm.File["file"]

			// store cached file ids.
			if len(fileHeaders) > 0 {
				fileHeader := fileHeaders[0]
				filename := fileHeader.Filename
				fileType := fileHeader.Header["Content-Type"][0]
				file, _ := fileHeader.Open()
				content, _ := ioutil.ReadAll(file)
				fileId := cFile.NewCacheFile(filename, fileType, content, 1)
				agentResp.SetParam("fileId", fileId)
				agentResp.SetRespCode(respcode.RC_GENERAL_SUCC)
			} else {
				agentResp.SetRespCode(respcode.RC_GENERAL_APP_ERR)
				agentResp.SetParam("error", "Dummy input file.")
			}

		}
	}
	respByte, _ := json.Marshal(agentResp)
	w.Write(respByte)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["fileId"]
	name, contentType, content, exists := cFile.GetCacheFile(fileId, false)

	if !exists {
		log.Debugf("load file %s request to file node", fileId)
		client, _ := discover.GetClient(action.Action_LoadFile)

		if client == nil {
			respBytes, _ := json.Marshal(respcode.RC_SERVICE_UNAVAILABLE)
			w.Write(respBytes)
			return
		}

		params := map[string]interface{}{"fileId":fileId}
		request := &message.Request{Action:action.Action_LoadFile, Params:params}
		cachedReq, err := message.CacheReq(request)

		clientResp, err := client.Serve(context.Background(), cachedReq)
		log.Debug("clientResp:", clientResp)
		if err == nil {
			respByte, _ := message.GetCacheResp(clientResp.ResponseKey)
			log.Debug(string(respByte))
			respObj := &message.Response{}
			err := json.Unmarshal(respByte, respObj)
			if err == nil {
				fileId, _ := respObj.Params["fileId"].(string)
				name, contentType, content, exists = cFile.GetCacheFile(fileId, false)
			} else {
				http.NotFoundHandler()
			}
		}
	} else {
		log.Debug("exists in cache")
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	w.Header().Add("Content-type", contentType)
	w.Write(content)
}