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

type EmailConf struct {
	Type      string
	Host      string
	Port      int
	EnableSsl bool
	UserName  string
	Password  string
	CronTime  string
}
