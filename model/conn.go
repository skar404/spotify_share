package model

import "gopkg.in/mgo.v2"

const dbName = "bot_db"

type Conn struct {
	DB *mgo.Session
}

func (c *Conn) Clone() (*mgo.Session, *mgo.Database) {
	conn := c.DB.Clone()
	return conn, conn.DB(dbName)
}
