/*
 * Revision History:
 *     Initial: 2018/05/21        Chen Yanchen
 */

package model

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/JonSnow47/Gym-management-system/db"
	"github.com/JonSnow47/Gym-management-system/util"
)

const collectionAdmin = "admin"

type adminServiceProvide struct{}

var AdminService *adminServiceProvide

func conAdmin() db.Connection {
	con := db.Connect(collectionAdmin)
	con.C.EnsureIndex(mgo.Index{
		Key:        []string{"_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})
	return con
}

type Admin struct {
	Id      bson.ObjectId `bson:"_id"`
	Name    string        `bson:"Name"`
	Pwd     string        `bson:"Pwd"`
	Created time.Time     `bson:"Created"`
}

// 新建管理员
func (*adminServiceProvide) New(name, pwd string) (string, error) {
	con := conAdmin()
	defer con.S.Close()

	if len(pwd) < 6 || len(pwd) > 16 {
		return "", errors.New("Password length error.")
	}

	data, err := util.GenerateHash(pwd)
	if err != nil {
		return "", err
	}

	var admin = &Admin{
		Id:      bson.NewObjectId(),
		Name:    name,
		Pwd:     string(data),
		Created: time.Now(),
	}

	err = con.C.Insert(admin)
	if err != nil {
		return "", err
	}

	return admin.Id.Hex(), nil
}

// Login
func (*adminServiceProvide) Login(name, pwd string) (bool, error) {
	con := conAdmin()
	defer con.D.Session.Close()

	var admin Admin

	err := con.C.Find(bson.M{"Name": name}).One(&admin)
	if err != nil {
		return false, err
	}

	if util.CompareHash([]byte(admin.Pwd), pwd) {
		return true, nil
	}
	return false, nil
}
