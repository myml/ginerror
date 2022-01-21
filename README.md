# ginerror

--
import "github.com/myml/ginerror"

## Usage

### func ErrorHandle

```go
func ErrorHandle(logger *zap.Logger, customHandle func(c *gin.Context, err error) bool) gin.HandlerFunc
```

ErrorHandle 错误处理中间件 统一处理参数错误，数据库错误和自定义错误

生成错误的 uuid，便于前后端跟踪错误

使用 Error()包裹错误后，可打印错误堆栈

参数错误返回 code 400，并附带错误的字段和校验器 例子{"code": 400,"meesage": "validate error","fields":
{"limit":"min"}}

数据库错误，当查询不到数据，返回 code 404

自定义错误，可通过 customHandle 自行处理，如非自定义错误，则返回 false

其他错误返回 500，并在日志打印 InternalServerError 日志

### func Errorf

```go
func Errorf(c *gin.Context, format string, err error) bool
```

Errorf 判断 err 是否为空 如果 err 为空返回 false，
如果 err 不为空返回 true，并执行记录堆栈，c.Error(fmt.Errorf(format, err)); c.Abort()。

可使用 errors.As(err,ErrorStack) 取出堆栈信息

类似 AbortWithError，但不设置 status，便于在错误处理中间件中自定义 http status code

### func RegisterTagName

```go
func RegisterTagName()
```

RegisterTagName 注册一个 tag name func 到 gin 的验证引擎 用于自动从 tag(优先级：uri,form,json)生成验证错误的字段名
使用 tag 做为错误字段名，有利于前端自动化提示错误字段 只支持新版本 gin 使用的 validator/v10

### type ErrorStack

```go
type ErrorStack struct {
}
```

ErrorStack 用于记录错误的堆栈信息

#### func (ErrorStack) Error

```go
func (es ErrorStack) Error() string
```

#### func (\*ErrorStack) Unwrap

```go
func (es *ErrorStack) Unwrap() error
```
