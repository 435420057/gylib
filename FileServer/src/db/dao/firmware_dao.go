package dao

import (
	"db"
	"db/entity"
	"gopkg.in/mgo.v2/bson"
	"github.com/kyugao/go-logger/logger"
)

func IsLatestFirmwareVersion(version string) (isLatest bool, latestVersion string, fileId string, err error) {
	collection := db.GetCollection(entity.ColFirmwareVersion)
	firmwareVersion := &entity.FirmwareVersion{}
	err = collection.Find(bson.M{"latest":true}).One(firmwareVersion)
	logger.Infof("find latest firmware version %v, with err %v", firmwareVersion, err)
	if err == nil {
		isLatest = (firmwareVersion.Version == version)
	}
	if !isLatest && firmwareVersion.FileId.Valid() {
		latestVersion = firmwareVersion.Version
		fileId = firmwareVersion.FileId.Hex()
	}
	return
}