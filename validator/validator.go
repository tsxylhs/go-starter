package validator

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/tsxylhs/go-starter/log"
)

var (
	validate *validator.Validate
)

func BindAndValid(c *gin.Context, domain interface{}) (int, string) {

	// if err := c.Bind(domain); err != nil {
	// 	log.Logger.Logger.Error("c.Bind is error", zap.Error(err))
	// 	return http.StatusBadRequest, 200
	// }
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
	}
	err = validate.Struct(domain)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Logger.Logger.Info(err.Translate(trans))

			return http.StatusInternalServerError, err.Translate(trans)
		}
	}

	return http.StatusOK, ""

}

func init() {
	validate = validator.New()
}
