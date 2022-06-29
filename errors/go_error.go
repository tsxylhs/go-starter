package errors

type ErrComponent string

//定义层级
const (
	ErrService ErrComponent = "service"
	ErrorRepo  ErrComponent = "repository"
	ErrLib     ErrComponent = "library"
	ErrorCgo   ErrComponent = "Cgo"
	ErrorDb    ErrComponent = "Db"
)

// 响应类型

type ResponseErrType string

const (
	BadRequest    ResponseErrType = "BadRequest"
	Forbiddens    ResponseErrType = "Forbidden"
	NotFounds     ResponseErrType = "NotFound"
	AlreadyExists ResponseErrType = "AlreadyExists"
)

//code 定义
type Code string

const (
	ServiceCode Code = "500"
)

var G GoError

type GoError struct {
	error
	Code         string                 //错误码
	Data         map[string]interface{} //描述
	Causes       []error                //原因
	Cause        error
	Component    ErrComponent    //错误层级
	ResponseType ResponseErrType //响应类型
	Retryable    bool
}

type GError interface {
	error
	Code() string
	Data() map[string]interface{}
	Causes() []error
	Component() ErrComponent       //错误层级
	ResponseType() ResponseErrType //响应类型
	Retryable() bool
	SetRetryalbe() GError
}

//todo 添加更多类型
func (g *GoError) ResourceNotFound(id, message string, cause error) GoError {
	data := map[string]interface{}{"message": message, "id": id}
	return GoError{
		Code:         string(ServiceCode),
		Data:         data,
		Causes:       []error{cause},
		Component:    ErrService,
		ResponseType: NotFounds,
		Retryable:    false,
	}
}
