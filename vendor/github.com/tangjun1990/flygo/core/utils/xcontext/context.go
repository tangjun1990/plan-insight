package xcontext

import (
	"context"
	"github.com/gin-gonic/gin"
	"unsafe"
)

// ToContext 获取 gin.context 中 request.context()
func ToContext(c context.Context) (ctx context.Context) {
	switch c.(type) {
	case *gin.Context:
		ctx = c.(*gin.Context).Request.Context()

		uid64 := c.(*gin.Context).GetInt64("uid")
		if uid64 != 0 {
			ctx = context.WithValue(ctx, "uid", uid64)
			return ctx
		}

		uid := c.(*gin.Context).GetInt("uid")
		if uid != 0 {
			ctx = context.WithValue(ctx, "uid", uid)
			return ctx
		}

		uidStr := c.(*gin.Context).GetString("uid")
		if uidStr != "" {
			ctx = context.WithValue(ctx, "uid", uidStr)
			return ctx
		}

		return ctx
	default:
		return c
	}
}

type iface struct {
	itab, data uintptr
}

type valueCtx struct {
	context.Context
	key, val interface{}
}

func GetKeyValues(ctx context.Context) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[interface{}]interface{}) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil && valCtx.val != nil {
		m[valCtx.key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}
