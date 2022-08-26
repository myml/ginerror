package ginerror

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Errorf 判断err是否为空
// 如果err为空返回false，
// 如果err不为空返回true，并执行c.Error c.Abort
//
// 可使用errors.As(err,ErrorStack) 取出堆栈信息
//
// 类似AbortWithError，但不设置status，便于在错误处理中间件中自定义http status code
func Errorf(c *gin.Context, format string, err error) bool {
	if err == nil {
		return false
	}
	pc, filename, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	c.Error(&ErrorStack{
		err:      fmt.Errorf(format, err),
		function: f.Name(),
		filename: fmt.Sprintf("%s:%d", filename, line),
	})
	c.Abort()
	return true
}

// ErrorHandle 错误处理中间件
// 统一处理参数错误，数据库错误和自定义错误
//
// 生成错误的uuid，便于前后端跟踪错误
//
// 使用Error()包裹错误后，可打印错误堆栈
//
// 参数错误返回code 400，并附带错误的字段和校验器 例子 {"code": 400,"meesage": "validate error","fields": {"limit":"min"}}
//
// 数据库错误，当查询不到数据，返回code 404
//
// 自定义错误，可通过customHandle自行处理，如非自定义错误，则返回false
//
// 其他错误返回500，并在日志打印InternalServerError日志
func ErrorHandle(logger *zap.Logger, customHandle func(c *gin.Context, err error) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		// 自定义处理函数，一般用于业务错误
		if customHandle != nil {
			if customHandle(c, err) {
				return
			}
		}

		requestID := uuid.New()

		var es ErrorStack
		if errors.As(err.Err, &es) {
			logger.Error("ErrorHandle",
				zap.String("requestID", requestID.String()),
				zap.String("requestURI", c.Request.RequestURI),
				zap.Error(err),
				zap.String("function", es.function),
				zap.String("filename", es.filename),
			)
		} else {
			logger.Error("ErrorHandle",
				zap.String("requestID", requestID.String()),
				zap.String("requestURI", c.Request.RequestURI),
				zap.Error(err),
			)
		}
		if c.Writer.Status() != http.StatusOK {
			return
		}

		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":      http.StatusNotFound,
				"meesage":   "resource not found",
				"requestID": requestID,
			})
			return
		}
		var validationErrors validator.ValidationErrors
		if errors.As(err.Err, &validationErrors) {
			fields := map[string]string{}
			for i := range validationErrors {
				fields[validationErrors[i].Field()] = validationErrors[i].Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"code":      http.StatusBadRequest,
				"meesage":   "validate error",
				"fields":    fields,
				"requestID": requestID,
			})
			return
		}
		// 未知错误再打印一次
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":      http.StatusInternalServerError,
			"message":   "unknown error",
			"requestID": requestID,
		})
	}
}
