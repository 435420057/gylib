package dao

import (
	"db"
	"db/entity"
	"gopkg.in/mgo.v2/bson"
)

func ValidAppVersion(platform, v, s, i, c string) (valid bool) {
	collection := db.GetCollection(entity.ColAppVersion)
	appVersion := &entity.AppVersion{}
	err := collection.Find(bson.M{"platform":platform, "version":v, "sub": s, "iteration":i, "channel": c}).One(appVersion)
	if err == nil {
		valid = appVersion.Valid
	} else {
		valid = false
	}
	return
}