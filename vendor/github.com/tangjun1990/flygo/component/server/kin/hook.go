package kin

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/tangjun1990/flygo/core/ktrace"
	"github.com/tangjun1990/flygo/core/utils/xcontext"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.uber.org/zap"

	"github.com/tangjun1990/flygo/core/kapp"
	"github.com/tangjun1990/flygo/core/klog"
	"github.com/tangjun1990/flygo/core/kmetric"
)

// 拦截器

const FGUID = "fg_uid"
const FGRID = "X-Request-ID"
const FGCID = "fg_company_id"

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func extractAPP(ctx *gin.Context) string {
	return ctx.Request.Header.Get("app")
}

type resWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (g *resWriter) Write(data []byte) (int, error) {
	n, e := g.body.Write(data)
	if e != nil {
		return n, e
	}
	return g.ResponseWriter.Write(data)
}

func (g *resWriter) WriteString(s string) (int, error) {
	n, e := g.body.WriteString(s)
	if e != nil {
		return n, e
	}
	return g.ResponseWriter.WriteString(s)
}

func defaultServerHook(logger *klog.Component, config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var beg = time.Now()
		var rw *resWriter

		b, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		rw = &resWriter{c.Writer, &bytes.Buffer{}}
		c.Writer = rw

		loggerKeys := kapp.OkLogExtraKeys()
		var fields = make([]klog.Field, 0, 20+len(loggerKeys))
		var brokenPipe bool
		defer func() {
			ctx := xcontext.ToContext(c)
			cost := time.Since(beg)

			agent := ""
			if values, ok := c.Request.Header["User-Agent"]; ok && len(values) > 0 {
				agent = values[0]
			}

			fields = append(fields,
				klog.FieldServerKind(),
				klog.FieldDuration(cost), // 耗时
				klog.FieldNetPeerIp(strings.Split(c.Request.RemoteAddr, ":")[0]), // 请求方式
				klog.FieldHttpMethod(c.Request.Method),                           // 请求方式
				klog.FieldHttpHost(c.Request.Host),                               // 实际请求域名
				klog.FieldHttpPath(c.Request.URL.Path),                           // 地址
				klog.FieldHttpTarget(c.Request.RequestURI),                       // target
				klog.FieldHttpUserAgent(agent),                                   // user agent
			)

			if values, ok := c.Request.Header["X-Forwarded-For"]; ok && len(values) > 0 {
				if addresses := strings.SplitN(values[0], ",", 2); len(addresses) > 0 {
					fields = append(fields,
						klog.FieldHttpClientIp(addresses[0]),
					)
				}
			}

			// 请求参数
			if config.HookReq {
				for k, v := range c.Request.Header {
					fields = append(fields, klog.Any("http.request.header."+strings.ReplaceAll(strings.ToLower(k), "-", "_"), strings.Join(v, ",")))
				}
				if b != nil {
					// 请求过长，则截断
					bodyStr := string(b)
					size := len(bodyStr)
					if size > config.MaxReqContentSize {
						bodyStr = bodyStr[:config.MaxReqContentSize] + " ..."
					}
					fields = append(fields, klog.String("http.request.body", bodyStr))
				}
			}

			// 请求返回
			if config.HookRsp {
				// 返回过长，则截断
				bodyStr := rw.body.String()
				size := len(bodyStr)
				if size > config.MaxResContentSize {
					bodyStr = bodyStr[:config.MaxResContentSize] + " ..."
				}

				for k, v := range c.Writer.Header() {
					fields = append(fields, klog.Any("http.response.header."+strings.ReplaceAll(strings.ToLower(k), "-", "_"), strings.Join(v, ",")))
				}

				fields = append(fields,
					klog.Int("http.response_content_length", size),
					klog.String("http.response.body", bodyStr),
				)
			}

			for _, key := range loggerKeys {
				if value := getContextValue(key, c); value != "" {
					fields = append(fields, klog.FieldCustomKeyValue(key, value))
				}
			}

			// 慢接口日志
			if config.SlowLogThreshold > time.Duration(0) && config.SlowLogThreshold < cost {
				logger.WithCtx(ctx).Warn("slow", fields...)
			}

			// 错误日志
			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				message := "panic"
				switch v := rec.(type) {
				case error:
					message = v.Error()
				case string:
					message = v
				}

				if brokenPipe {
					_ = c.Error(rec.(error))
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}

				stackInfo := stack(3)

				fields = append(fields,
					zap.ByteString("stack", stackInfo),
					klog.FieldEvent("panic"),
					klog.FieldErrAny(rec),
					klog.FieldHttpStatusCode(c.Writer.Status()),
				)
				logger.WithCtx(ctx).Error(message, fields...)
				return
			}

			// 正常日志
			fields = append(fields,
				klog.FieldHttpStatusCode(c.Writer.Status()),
			)

			if c.Errors.ByType(gin.ErrorTypePrivate).String() != "" {
				fields = append(fields,
					klog.FieldErrAny(c.Errors.ByType(gin.ErrorTypePrivate).String()),
				)
				logger.WithCtx(ctx).Error(c.Errors.ByType(gin.ErrorTypePrivate).String(), fields...)
			} else {
				logger.WithCtx(ctx).Info("access", fields...)
			}

		}()

		c.Next()
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func metricServerHook() gin.HandlerFunc {
	return func(c *gin.Context) {
		beg := time.Now()
		c.Next()
		kmetric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), kmetric.TypeHTTP, c.Request.Method+"."+c.FullPath(), extractAPP(c))
		kmetric.ServerHandleCounter.Inc(kmetric.TypeHTTP, c.Request.Method+"."+c.FullPath(), extractAPP(c), http.StatusText(c.Writer.Status()))
	}
}

func getPeerIP(addr string) string {
	addSlice := strings.Split(addr, ":")
	if len(addSlice) > 1 {
		return addSlice[0]
	}
	return ""
}

func getContextValue(key string, c *gin.Context) string {
	if key == "" {
		return ""
	}
	val := cast.ToString(c.Request.Context().Value(key))
	if val == "" {
		return c.GetHeader(key)
	}
	return val
}

func traceServerHook() gin.HandlerFunc {
	tracer := ktrace.NewTracer(trace.SpanKindServer)
	attrs := []attribute.KeyValue{
		semconv.RPCSystemKey.String("http"),
	}
	return func(c *gin.Context) {
		b, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		reqBody := ""
		if b != nil {
			reqBody = string(b)
		}

		authHeader := c.GetHeader("Febook-Gateway-User")
		if len(authHeader) > 0 {
			authHeaderByte, err := base64.StdEncoding.DecodeString(authHeader)
			if err == nil {
				authHeader = string(authHeaderByte)
			}
		}

		ktrace.CompatibleExtractHTTPTraceID(c.Request.Header)

		ctx, span := tracer.Start(c.Request.Context(), c.Request.Method+"."+c.FullPath(), propagation.HeaderCarrier(c.Request.Header), trace.WithAttributes(attrs...))
		span.SetAttributes(
			semconv.HTTPURLKey.String(c.Request.URL.String()),
			semconv.HTTPTargetKey.String(c.Request.URL.Path),
			semconv.HTTPMethodKey.String(c.Request.Method),
			semconv.HTTPUserAgentKey.String(c.Request.UserAgent()),
			semconv.HTTPClientIPKey.String(c.ClientIP()),
			ktrace.CustomTag("http.full_path", c.FullPath()),
			ktrace.CustomTag("http.request.body", reqBody),
			ktrace.CustomTag("http.request.header.Febook-Gateway-User", authHeader),
		)
		c.Request = c.Request.WithContext(ctx)
		c.Header(FGRID, span.SpanContext().TraceID().String())
		var rw *resWriter

		rw = &resWriter{c.Writer, &bytes.Buffer{}}
		c.Writer = rw

		c.Next()
		uid, ok := c.Request.Context().Value(FGUID).(int)
		uidstr := ""
		if ok {
			uidstr = cast.ToString(uid)
		}

		cid, ok := c.Request.Context().Value(FGCID).(int)
		cidstr := ""
		if ok {
			cidstr = cast.ToString(cid)
		}
		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int64(int64(c.Writer.Status())),
			ktrace.CustomTag("feige.uid", uidstr),
			ktrace.CustomTag("feige.companyid", cidstr),
			ktrace.CustomTag("feige.traceid", c.GetHeader(FGRID)),
			ktrace.CustomTag("http.response.body", rw.body.String()),
		)
		span.End()
	}
}
