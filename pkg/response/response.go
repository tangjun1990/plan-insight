package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 接口返回
type Response struct {
	ctx       *gin.Context
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	RequestID string      `json:"requestID"`
	Data      interface{} `json:"data"`
	Status    bool        `json:"status"`
}

func WithCtx(ctx *gin.Context) *Response {
	r := &Response{}
	r.ctx = ctx
	r.RequestID = ctx.Writer.Header().Get("x-trace-id")
	return r
}

func (r *Response) Success(d interface{}) {
	r.Msg = "success"
	r.Status = true
	r.Data = d
	r.ctx.JSON(http.StatusOK, r)
}

func (r *Response) Error(code int, msg ...string) {
	r.Code = code
	r.Status = false
	if len(msg) != 0 {
		r.Msg = msg[0]
	} else {
		r.Msg = message(code)
	}
	r.ctx.JSON(http.StatusOK, r)
}
