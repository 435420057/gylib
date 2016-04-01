package db

import (
	"external/mongo"
	log "github.com/kyugao/go-logger/logger"
	"gopkg.in/mgo.v2"
	"github.com/stvp/go-toml-config"
)

const db_config_path = "./conf/db.conf"

var (
	dialContext *mongo.DialContext
	dbConfig *config.ConfigSet

	dbUrl string
	DBName string
	dbConnPoolSize int
)

func init() {
	loadConfig()
	mongodbContext, err := mongo.Dial(dbUrl, dbConnPoolSize)
	if err != nil {
		log.Debug("Could not establish connection with db:", err)
	} else {
		log.Debugf("Connected to db %s.", dbUrl)
		dialContext = mongodbContext
	}
}

func loadConfig() {
	dbConfig = config.NewConfigSet("dbConfig", config.ExitOnError)
	dbConfig.StringVar(&dbUrl, "db_url", "mongodb://glareme.cn:27017")
	dbConfig.StringVar(&DBName, "db_name", "file_schema")
	dbConfig.IntVar(&dbConnPoolSize, "db_conn_pool_size", 5)
	err := dbConfig.Parse(db_config_path)
	if err != nil {
		log.Warnf("load dbconfig error, %v", err)
	} else {
		log.Info("loaded dbconfig")
	}
}

func getSession() *mongo.Session {
	context := dialContext.Ref()
	return context
}

func returnSession(session *mongo.Session) {
	dialContext.UnRef(session)
}

func GetCollection(name string) *mgo.Collection {
	session := getSession()
	defer returnSession(session)
	return session.Session.DB(DBName).C(name)
}

func GetGridFS(category string) *mgo.GridFS{
	session := getSession()
	defer returnSession(session)
	return session.DB(DBName).GridFS(category)
}