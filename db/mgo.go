package db

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/olebedev/config"
	"gopkg.in/mgo.v2"
)

type (
	Connect struct {
		Session *mgo.Session
		DbName  string
	}

	Collection struct {
		MemUsageCollection *mgo.Collection
		CpuUsageCollection *mgo.Collection
	}
)

func New(cfg *config.Config) *Connect {

	mInfo, err := getInfo(cfg)

	if err != nil {
		log.Fatalf("mgo DiaInfo get Error : ", err)

	}

	mSession, err := mgo.DialWithInfo(mInfo)

	if err != nil {
		panic(err)
	}

	defer mSession.Close()

	mSession.SetSafe(&mgo.Safe{})

	return &Connect{
		Session: mSession.Copy(),
		DbName:  mInfo.Database,
	}

}

func getInfo(cfg *config.Config) (*mgo.DialInfo, error) {

	mHost, err := cfg.String("development.database.host")
	if err != nil {
		log.Fatalf("get config host error : ", err)
		return nil, err
	}

	mPort, err := cfg.String("development.database.port")

	if err != nil {
		log.Fatalf("get config port error : ", err)
		return nil, err
	}

	mUser, err := cfg.String("development.database.username")

	if err != nil {
		log.Fatalf("get config username error : ", err)
		return nil, err
	}

	mPwd, err := cfg.String("development.database.password")

	if err != nil {
		log.Fatalf("get config password error : ", err)
		return nil, err
	}

	mdbName, err := cfg.String("development.database.dbname")

	if err != nil {
		log.Fatalf("get config dbname error : ", err)
		return nil, err
	}

	return &mgo.DialInfo{
		Addrs:    []string{strings.Join([]string{mHost, mPort}, ":")},
		Timeout:  60 * time.Second,
		Database: mdbName,
		Username: mUser,
		Password: mPwd,
	}, nil

}

func GetColl(c *Connect) *Collection {

	return &Collection{
		MemUsageCollection: c.Session.DB(c.DbName).C("memstats"),
		CpuUsageCollection: c.Session.DB(c.DbName).C("cpustats"),
	}

}

func Save(coll *mgo.Collection, doc interface{}) error {

	err := coll.Insert(doc)

	if err != nil {
		return err
	}

	return nil

}
