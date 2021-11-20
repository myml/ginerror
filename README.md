# ginerror
--
    import "github.com/myml/ginerror"


## Usage

#### func  ErrorHandle

```go
func ErrorHandle(logger *zap.Logger, customHandle func(c *gin.Context, err error) bool) gin.HandlerFunc
```
ErrorHandle 错误处理中间件 统一处理参数错误，数据库错误和自定义错误

生成错误的uuid，便于前后端跟踪错误

使用Error()包裹错误后，可打印错误堆栈

参数错误返回code 400，并附带错误的字段和校验器 例子{"code": 400,"meesage": "validate error","fields":
{"limit":"min"}}

数据库错误，当查询不到数据，返回code 404

自定义错误，可通过customHandle自行处理，如非自定义错误，则返回false

其他错误返回500，并在日志打印InternalServerError日志

#### func  Errorf

```go
func Errorf(c *gin.Context, format string, err error) bool
```
Errorf 判断err是否为空 如果err为空返回false，
如果err不为空返回true，并执行记录堆栈，c.Error(fmt.Errorf(format, err)); c.Abort()。

可使用errors.As(err,ErrorStack) 取出堆栈信息

类似AbortWithError，但不设置status，便于在错误处理中间件中自定义http status code

#### func  RegisterTagName

```go
func RegisterTagName()
```
RegisterTagName 注册一个tag name func到gin的验证引擎 用于自动从tag(优先级：uri,form,json)生成验证错误的字段名
使用tag做为错误字段名，有利于前端自动化提示错误字段 只支持新版本gin使用的 validator/v10

#### type ErrorStack

```go
type ErrorStack struct {
}
```

ErrorStack 用于记录错误的堆栈信息

#### func (ErrorStack) Error

```go
func (es ErrorStack) Error() string
```

#### func (*ErrorStack) Unwrap

```go
func (es *ErrorStack) Unwrap() error
```
