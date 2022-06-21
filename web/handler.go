package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
	starterModle "github.com/tsxylhs/go-starter"
	"github.com/tsxylhs/go-starter/log"
)

type Handler struct {
}

func (handler *Handler) Register(router *gin.Engine) {
	log.Logger.Fatal("not implemented")
}

func (handler *Handler) UID(c *gin.Context) int64 {
	user, ok := c.Get(starterModle.UserKey)
	if !ok || user == nil {
		return 0
	}
	return user.(starterModle.IdInf).GetId()
}

func (handler *Handler) User(c *gin.Context) interface{} {
	user, _ := c.Get(starterModle.UserKey)
	return user
}

func (handler *Handler) Bind(c *gin.Context, domain interface{}) (err error) {
	err = c.ShouldBind(domain)
	if err != nil {
		//return error(&errors.SimpleBizError{Code: errors.Common_InvalidParams, Msg: err.Error()})
	}
	return nil
}

// func (handler *Handler) BindAndValidate(c *gin.Context, domain interface{}, ruleSetName string) (err error) {
// 	err = c.ShouldBind(domain)
// 	if err != nil {
// 		return error(&errors.SimpleBizError{Code: errors.Common_InvalidParams, Msg: err.Error()})
// 	}

// 	return validator.ValidateStruct(domain, ruleSetName)
// }

func (handler *Handler) Int64Param(c *gin.Context, key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}
