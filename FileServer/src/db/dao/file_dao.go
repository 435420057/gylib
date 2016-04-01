package dao
import (
	"gopkg.in/mgo.v2/bson"
	"db"
	"gopkg.in/mgo.v2"
	"errors"
	log "github.com/kyugao/go-logger/logger"
)

const category = "gridfs"

func LoadFile(fileId string) (name string, contentType string, content []byte, err error) {

	if (!bson.IsObjectIdHex(fileId)) {
		err = errors.New("invalid id hex string.")
		return
	}

	gridFS := db.GetGridFS(category)
	gridFile, err := gridFS.OpenId(bson.ObjectIdHex(fileId))

	log.Debugf("find grid file %v, err %v:", gridFile, err)
	if err == mgo.ErrNotFound {
		return
	}
	name = gridFile.Name()
	log.Debugf("read filename, ", name)
	contentType = gridFile.ContentType()
	content = make([]byte, gridFile.Size())
	gridFile.Read(content)
	return
}

func SaveFile(name string, contentType string, content []byte) (collection string, fileId string, dbName string, err error) {
	gridFS := db.GetGridFS(category)
	file, err := gridFS.Create(name)
	if err != nil {
		return
	} else {
		fileId = file.Id().(bson.ObjectId).Hex()
		collection = category
		dbName = db.DBName
		file.SetContentType(contentType)
		num, err := file.Write(content)
		log.Debug(num, ":", err)
		file.Close()
	}
	return
}

func DeleteFile(fileId string) (err error) {
	gridFS := db.GetGridFS(category)
	if (!bson.IsObjectIdHex(fileId)) {
		err = errors.New("invalid id hex string.")
	} else {
		err = gridFS.RemoveId(bson.ObjectIdHex(fileId))
	}
	return
}