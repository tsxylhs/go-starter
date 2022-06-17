package code

import "xorm.io/xorm"

type DBListener interface {
	SetDB(db *xorm.Engine)
	GetDbName() string
	SetDbName(string)
	DbEnabled() bool
}

type IModule interface {
	DBListener
	GetName() string
	GetTableName() string
}
