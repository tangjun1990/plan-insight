package ktrace

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tangjun1990/flygo/core/utils/xstring"

	"google.golang.org/grpc/metadata"
)

// CompatibleExtractHTTPTraceID ...
// Deprecated 该方法会在v1.2.0移除
func CompatibleExtractHTTPTraceID(header http.Header) {
	xTraceID := header.Get("X-Request-ID")
	if xTraceID != "" {
		header.Set("Traceparent", CompatibleParse(xTraceID))
	}
}

// CompatibleExtractGrpcTraceID ...
// Deprecated 该方法会在v1.2.0移除
func CompatibleExtractGrpcTraceID(header metadata.MD) {
	xTraceID := header.Get("x-request-id")
	if len(xTraceID) > 0 {
		header.Set("Traceparent", CompatibleParse(xTraceID[0]))
	}
}

// CompatibleParse ...
// opentrace: 18af9db18a77f4b7:18af9db18a77f4b7:0000000000000000:0
// opentelemetry: 00-18af9db18a77f4b718af9db18a77f4b7-18af9db18a77f4b7-00
// https://www.w3.org/TR/trace-context/
func CompatibleParse(traceID string) string {
	if len(traceID) == 32 {
		return RequestIDToParent(traceID)
	}
	traceArr := strings.Split(traceID, ":")
	if len(traceArr) == 4 {
		return "00-" + traceArr[0] + traceArr[1] + "-" + traceArr[1] + "-0" + traceArr[3]
	}
	return ""
}

// RequestIDToParent 通过requestID,生成完整的 traceparent
func RequestIDToParent(traceID string) string {
	if len(traceID) != 32 {
		return ""
	}

	// version-traceID-spanID-traceFlags
	return fmt.Sprintf("00-%s-%s-01", traceID, xstring.RandByteStr(8))
}
