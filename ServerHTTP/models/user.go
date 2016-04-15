package models

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// Users nombre de la colección de usuarios
	Users = "users"
	// UsernameKey es el key que usara index en la colección de users
	UsernameKey = "username"
	// PasswordKey es el key para la contraseña del usuario
	PasswordKey = "password"
)

// User es la estructura de un usuario
type User struct {
	ID       bson.ObjectId
	Username string
	Password string
}

// UserModel es la estructura del modelo de usuarios
type UserModel struct{}

func (model *UserModel) insertIndexes() {
	db := DB{}
	err := db.Conn()

	if err == nil {
		collection := db.Collection(Users)

		hasIndex := false
		indexes := []mgo.Index{}

		indexes, err = collection.Indexes()
		if err == nil {
			for i := 0; i < len(indexes); i++ {
				index := indexes[i]
				keys := index.Key
				for j := 0; j < len(keys); j++ {
					key := keys[j]
					if key == UsernameKey {
						hasIndex = true
						break
					}
				}
			}

			if !hasIndex {
				index := mgo.Index{
					Key:    []string{UsernameKey},
					Unique: true,
				}

				hasIndex = true
				collection.EnsureIndex(index)
			}

			fmt.Println("users.insertIndexes: ", hasIndex, err)
		} else {
			fmt.Println("users.Indexes: ", err)
		}

	} else {
		fmt.Println("insertIndexes: ", err)
	}
}

// Find es el metodo para obtener un conjunto de usuarios
func (model *UserModel) Find() ([]User, error) {
	docs := []User{}
	db := DB{}
	err := db.Conn()

	if err == nil {
		collection := db.Collection(Users)
		err = collection.Find(bson.M{}).All(&docs)
		if err == nil {
			fmt.Println("Docs: ", docs)
		} else {
			fmt.Println("Find.All: ", err)
		}
	} else {
		fmt.Println("db.Conn: ", err)
	}

	return docs, err
}

// FindOne es el metodo para obtener un usuario
func (model *UserModel) FindOne(where bson.M) (User, error) {
	doc := User{}
	db := DB{}
	err := db.Conn()

	if err == nil {
		collection := db.Collection(Users)
		err = collection.Find(where).One(&doc)
		if err == nil {
			fmt.Println("Doc: ", doc)
		} else {
			fmt.Println("Find.One: ", err)
		}
	} else {
		fmt.Println("db.Conn: ", err)
	}

	return doc, err
}

// Insert es el metodo para crear un usuario
func (model *UserModel) Insert(doc User) (User, error) {
	db := DB{}
	err := db.Conn()
	if err == nil {
		hash := sha512.New()
		hash.Write([]byte(doc.Password))
		pwd := base64.URLEncoding.EncodeToString(hash.Sum(nil))
		doc.Password = pwd

		collection := db.Collection(Users)
		doc.ID = bson.NewObjectId()
		err = collection.Insert(&doc)
	}
	return doc, err
}

// ComparePassword verifica que el usuario y contraseña sean correctos
func (user *User) ComparePassword() bool {
	isMatch := false
	if user.Password != "" {
		hash := sha512.New()
		hash.Write([]byte(user.Password))
		pwd := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		db := DB{}
		err := db.Conn()
		if err == nil {
			defer db.Disconn()
			doc := User{}
			collection := db.Collection(Users)
			where := bson.M{UsernameKey: user.Username, PasswordKey: pwd}
			err = collection.Find(where).One(&doc)
			if err == nil && doc.ID != "" {
				isMatch = true
			}
		}
	}

	return isMatch
}
