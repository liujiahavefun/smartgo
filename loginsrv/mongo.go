package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id     bson.ObjectId `bson:"_id"`
	Name   string        `bson:"name"`
	Passwd string        `bson:"passwd"`
	Uid    string        `bson:"uid"`
}

const MONGO_URL = "127.0.0.1:27017"

var (
	mgoSession    *mgo.Session
	userDatabase  = "userdb"
	userColletion = "user"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(MONGO_URL)
		if err != nil {
			panic(err)
		}
	}

	//最大连接池默认为4096
	return mgoSession.Clone()
}

/**
 * 公共方法，获取collection对象
 */
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getMongoSession()
	defer session.Close()
	c := session.DB(userDatabase).C(collection)
	return s(c)
}

/**
 * 获取一条记录通过objectId
 */
func GetPersonById(id string) (*User, error) {
	objId := bson.ObjectIdHex(id)
	userObj := new(User)
	query := func(c *mgo.Collection) error {
		return c.FindId(objId).One(&userObj)
	}

	err := witchCollection(userColletion, query)
	return userObj, err
}

/**
 * 获取一条记录通过name
 */
func GetPersonByName(name string) (*User, error) {
	userObj := new(User)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"name": name}).One(&userObj)
	}

	err := witchCollection(userColletion, query)
	return userObj, err
}

/**
 * 获取一条记录通过uid
 */
func GetPersonByUid(uid string) (*User, error) {
	userObj := new(User)
	query := func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&userObj)
	}

	err := witchCollection(userColletion, query)
	return userObj, err
}
