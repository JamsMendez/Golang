package models

import "gopkg.in/mgo.v2"

// Configuración de la base de datos
const (
	Host     = "127.0.0.1:27017"
	Database = "server_go"
)

// Model es la interface para los modelos
type Model interface {
	Find()
	FindOne()
	Insert()
	Update()
	Remove()
}

// DB es la estructura para la base de datos
type DB struct {
	Session  *mgo.Session
	Database *mgo.Database
}

// Conn realizar la conexion a la base de datos
func (db *DB) Conn() error {
	var err error
	db.Session, err = mgo.Dial(Host)
	if err == nil {
		db.Session.SetMode(mgo.Monotonic, true)
		db.Database = db.Session.DB(Database)
	}

	return err
}

// Disconn realizar la conexion a la base de datos
func (db *DB) Disconn() {
	db.Session.Close()
}

// Collection devuelve la collection correspondiente
func (db *DB) Collection(name string) *mgo.Collection {
	return db.Database.C(name)
}
