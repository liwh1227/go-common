package response

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

const (
	Code    = "code"
	Message = "message"
	Data    = "data"
	Detail  = "detail"
)

// Response 通用的返回结构
type Response struct {
	// 错误码
	Code ErrCode `json:"code"`
	// 消息
	Message string `json:"message"`
	// 	错误细节
	Detail string `json:"detail,omitempty"`
	// 返回的结构体
	Data interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	retH := gin.H{
		Code:    ERR_CODE_OK,
		Message: ERR_CODE_OK.String(),
	}
	if isNil(data) {
		ctx.JSON(http.StatusOK, retH)
		return
	}
	retH[Data] = data
	ctx.JSON(http.StatusOK, retH)
}

func FailResponse(ctx *gin.Context, err error) {
	retH := gin.H{
		Code:    getCodeFromCommonErr(err),
		Message: getCodeFromCommonErr(err).String(),
	}
	if len(err.Error()) != 0 {
		retH[Detail] = err.Error()
		ctx.JSON(http.StatusOK, retH)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, retH)
	ctx.Abort()
}

// 判断返回数据是否为空
func isNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}

func getCodeFromCommonErr(err error) ErrCode {
	val, ok := err.(*Error)
	if !ok {
		return ERR_CODE_INTERNAL_UNKNOWN_ERR_TYPE
	}

	return val.Code
}
