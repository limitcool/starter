package code

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"google.golang.org/grpc/status"
)

// ApiResponse 代表一个响应给客户端的消息结构，包括错误码，错误消息，响应数据
type ApiResponse struct {
	RequestId string      `json:"-"`              // 请求的唯一ID
	ErrorCode int         `json:"code"`           // 错误码，0表示无错误
	Message   string      `json:"msg"`            // 提示信息
	Data      interface{} `json:"data,omitempty"` // 响应数据，一般从这里前端从这个里面取出数据展示
}

// 正确返回并返回空数据
func NullData(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": nil, "msg": code.GetMsg(0)})
	ctx.AbortWithStatusJSON(http.StatusOK, ApiResponse{
		// RequestId: cast.ToString(ctx.Value(global.RequestIDKey)),
		ErrorCode: 0,
		Message:   GetMsg(Success),
		Data:      []int{},
	})
}

// 正确返回并返回数据
func Ok(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, ApiResponse{
		// RequestId: cast.ToString(ctx.Value(global.RequestIDKey)),
		ErrorCode: 0,
		Message:   GetMsg(Success),
		Data:      data,
	})
}

// 错误返回
func Err(ctx *gin.Context, ecode int, msg ...string) {
	message := GetMsg(ecode)
	if len(msg) > 0 {
		message = msg[0]
	}
	ctx.JSON(http.StatusBadRequest, ApiResponse{
		// RequestId: cast.ToString(ctx.Value(global.RequestIDKey)),
		ErrorCode: ecode,
		Message:   message,
		Data:      []int{},
	})

}

// 自动返回请求体
func AutoResponse(ctx *gin.Context, data interface{}, err error) {
	if err == nil {
		if data == nil {
			NullData(ctx)
		} else {
			Ok(ctx, data)
		}
	} else {
		//错误返回
		errcode := ErrorUnknown
		errmsg := "服务器开小差啦，稍后再来试一试"
		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			// logx.Debug("未定义的错误:", err)
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if IsCodeErr(cast.ToInt(grpcCode)) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = cast.ToInt(grpcCode)
					errmsg = gstatus.Message()
				}
			}
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ApiResponse{
			// RequestId: cast.ToString(r.Context().Value(global.RequestIDKey)),
			ErrorCode: errcode,
			Message:   errmsg,
			Data:      []int{},
		})
	}
}
