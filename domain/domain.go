package code

import (
	"github.com/gin-gonic/gin"
	"github.com/tsxylhs/go-starter/db"
	"github.com/tsxylhs/go-starter/errors"
	"xorm.io/xorm"
)

const (
	UserKey          = "user"
	UserIdKey        = "uid"
	UserFirstNameKey = "firstName"
	UserLastNameKey  = "lastName"
	UserEmailKey     = "email"
	UserNicknameKey  = "nickname"
	UserRoleKey      = "role"
	UserRightKey     = "rights"
	UserTypeKey      = "tp"
	UserOrgIdKey     = "orgId"
	UserGroupKey     = "group"
)

type IdInf interface {
	SetId(int64)
	GetId() int64
}

type ID struct {
	Id int64 `xorm:"pk BIGINT(20)" json:"id,string" form:"id"`
}

func (idb *ID) SetId(id int64) {
	idb.Id = id
}

func (idb *ID) GetId() int64 {
	return idb.Id
}

type Context map[string]interface{}

func (ctx *Context) Get(key string) interface{} {
	return (*ctx)[key]
}

func (ctx *Context) MustGet(key string) interface{} {
	v := (*ctx)[key]

	if v == nil {
		panic("key " + key + " not present in context")
	}
	return v
}

func (ctx *Context) Set(key string, value interface{}) {
	(*ctx)[key] = value
}

type IResult interface {
	IsOk() bool
	Err() errors.BizError
	SetError(err errors.BizError)
	Set(key string, value interface{})
}

type Result struct {
	Ok    bool                   `json:"ok"`
	Error errors.BizError        `json:"err,omitempty"`
	Data  interface{}            `json:"data,omitempty"`
	User  interface{}            `json:"user,omitempty"`
	Extra map[string]interface{} `json:"extra,omitempty"`
}

func (r *Result) IsOk() bool {
	return r.Ok
}

func (r *Result) FillUser(c *gin.Context) {
	r.User, _ = c.Get(UserKey)
}

func (r *Result) Set(key string, value interface{}) {
	if r.Extra == nil {
		r.Extra = map[string]interface{}{}
	}
	r.Extra[key] = value
}

func (r *Result) Err() errors.BizError {
	return r.Error
}

func (r *Result) SetError(err errors.BizError) {
	r.Error = err
}

func (r *Result) Failure(errs ...errors.BizError) *Result {
	r.Ok = false
	if len(errs) > 0 {
		r.Error = errs[0]
	}
	return r
}

func (r *Result) FailureWithData(data interface{}, err errors.BizError) *Result {
	r.Ok = false
	r.Error = err
	r.Data = data

	return r
}

func (r *Result) Success(ds ...interface{}) *Result {
	r.Ok = true
	if len(ds) > 0 {
		r.Data = ds[0]
	}
	return r
}

func NewResult(data interface{}) *Result {
	return &Result{Error: &errors.SimpleBizError{}, Data: data}
}

type FilterResult struct {
	Result
	Page *Page `json:"page,omitempty"`
}

type Page struct {
	P   int    `json:"p" form:"p"`
	Ps  int    `json:"ps" form:"ps"`
	Cnt int64  `json:"cnt" form:"cnt"`
	K   string `json:"k" form:"k"`
	Pc  int    `json:"pc" form:"pc"`
	Od  string `json:"od,omitempty" form:"od"`
}

func (page *Page) GetPage() *Page {
	return page
}

func (page *Page) GetPager(count int64) *Page {
	page.Cnt = count
	if page.P < 1 {
		page.P = 1
	}
	if page.Ps < 1 {
		page.Ps = db.DEFAULT_PAGE_SIZE
	}
	page.Pc = int(page.Cnt)/page.Ps + 1
	return page
}

func (page *Page) Skip() int {
	if page.Ps > 0 {
		return (page.P - 1) * page.Ps
	}

	return (page.P - 1) * db.DEFAULT_PAGE_SIZE
}

func (page *Page) Limit() int {
	if page.Ps > 0 {
		return page.Ps
	}

	return db.DEFAULT_PAGE_SIZE
}

type Filter interface {
	Apply(session *xorm.Session)
	GetPage() *Page
}

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
