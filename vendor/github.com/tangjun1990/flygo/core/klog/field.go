package klog

import (
	"context"
	"github.com/spf13/cast"
	"strings"
	"time"

	"go.uber.org/zap"
)

// FieldComponent 设置组件
func FieldComponent(value string) Field {
	return String("cpt", value)
}

// FieldComponentName 设置组件配置名
func FieldComponentName(value string) Field {
	return String("cpt_name", value)
}

// FieldApp 设置应用名
func FieldApp(value string) Field {
	return String("app", value)
}

// FieldServiceName 设置服务名称
func FieldServiceName(serviceName string) Field {
	return String("service.name", serviceName)
}

// FieldServiceTenantId serviceTenantId
func FieldServiceTenantId(serviceTenantId string) Field {
	return String("service.tenant.id", serviceTenantId)
}

// FieldDeploymentEnvironment deploymentEnvironment
func FieldDeploymentEnvironment(deploymentEnvironment string) Field {
	return String("deployment.environment", deploymentEnvironment)
}

// FieldCloudRegion cloudRegion
func FieldCloudRegion(cloudRegion string) Field {
	return String("cloud.region", cloudRegion)
}

// FieldK8sNodeIp k8sNodeIp
func FieldK8sNodeIp(k8sNodeIp string) Field {
	return String("k8s.node.ip", k8sNodeIp)
}

// FieldNetPeerIp netPeerIp
func FieldNetPeerIp(netPeerIp string) Field {
	return String("net.peer.ip", netPeerIp)
}

// FieldNetPeerName netPeerName
func FieldNetPeerName(netPeerName string) Field {
	return String("net.peer.name", netPeerName)
}

// FieldNetPeerPort netPeerPort
func FieldNetPeerPort(netPeerPort int) Field {
	return Int("net.peer.port", netPeerPort)
}

// FieldHttpClientIp httpClientIp
func FieldHttpClientIp(httpClientIp string) Field {
	return String("http.client_ip", httpClientIp)
}

// FieldNetHostIp netHostIp
func FieldNetHostIp(netHostIp string) Field {
	return String("net.host.ip", netHostIp)
}

// FieldClientIp client_ip
func FieldClientIp(ip string) Field {
	return FieldNetPeerIp(ip)
}

// FieldServerIp server_ip
func FieldServerIp(ip string) Field {
	return FieldNetHostIp(ip)
}

// FieldAddr 设置地址
func FieldAddr(value string) Field {
	return String("addr", value)
}

// FieldName ...
func FieldName(value string) Field {
	return String("name", value)
}

// FieldTableName 数据库名称
func FieldTableName(value string) Field {
	return FieldDBSQLTable(value)
}

// FieldDBSQLTable 数据库名称
func FieldDBSQLTable(value string) Field {
	return String("db.sql.table", value)
}

// FieldDBName 数据库名称
func FieldDBName(value string) Field {
	return String("db.name", value)
}

// FieldDBStatement ...
func FieldDBStatement(value string) Field {
	return String("db.statement", value)
}

// FieldDBResult ...
func FieldDBResult(value string) Field {
	return String("db.result", value)
}

// FieldType ... level 1
// 禁止在日志中记录 type 字段
func FieldType(value string) Field {
	return String("log_type", value)
}

// FieldKind ... level 2
func FieldKind(value string) Field {
	return String("kind", value)
}

func FieldClientKind() Field {
	return FieldKind("client")
}

func FieldServerKind() Field {
	return FieldKind("server")
}

func FieldProducerKind() Field {
	return FieldKind("producer")
}

func FieldConsumerKind() Field {
	return FieldKind("consumer")
}

// FieldCode 不推荐 未来会废弃 请使用语义更加明确的 FieldHttpCode
func FieldCode(value int32) Field {
	return FieldHttpStatusCode(int(value))
}

// FieldHttpStatusCode ...
func FieldHttpStatusCode(httpStatusCode int) Field {
	return Int("http.status_code", httpStatusCode)
}

// FieldHttpPath ...
func FieldHttpPath(httpPath string) Field {
	return String("http.path", httpPath)
}

// FieldHttpTarget ...
func FieldHttpTarget(httpTarget string) Field {
	return String("http.target", httpTarget)
}

// FieldHttpUserAgent ...
func FieldHttpUserAgent(httpUserAgent string) Field {
	return String("http.user_agent", httpUserAgent)
}

// FieldHttpHost ...
func FieldHttpHost(httpHost string) Field {
	return String("http.host", httpHost)
}

// FieldTid 设置链路id 不推荐 未来会废弃 请使用语义更加明确的 FieldTraceID
func FieldTid(value string) Field {
	return String("trace_id", value)
}

// FieldTraceID 设置链路id
func FieldTraceID(value string) Field {
	return String("trace_id", value)
}

// FieldUid 设置用户Id
func FieldUid(uid string) Field {
	return String("uid", uid)
}

// FieldCtxUid 从 context 中获取 uid
func FieldCtxUid(ctx context.Context) Field {
	uid, ok := ctx.Value("feige.uid").(int)
	uidstr := ""
	if ok {
		uidstr = cast.ToString(uid)
	}

	return FieldUid(uidstr)
}

// FieldSize ...
func FieldSize(value int32) Field {
	return Int32("size", value)
}

// FieldCost 耗时时间 ms
func FieldCost(value time.Duration) Field {
	return FieldDuration(value)
}

// FieldDuration 耗时时间 ms
func FieldDuration(duration time.Duration) Field {
	return zap.Float64("duration", float64(duration.Microseconds())/1000)
}

// FieldKey ...
func FieldKey(value string) Field {
	return String("key", value)
}

// FieldValue ...
func FieldValue(value string) Field {
	return String("value", value)
}

// FieldValueAny ...
func FieldValueAny(value interface{}) Field {
	return Any("value", value)
}

// FieldErrKind ...
func FieldErrKind(value string) Field {
	return String("errKind", value)
}

// FieldErr ...
func FieldErr(err error) Field {
	return zap.Error(err)
}

// FieldErrAny ...
func FieldErrAny(err interface{}) Field {
	return zap.Any("error", err)
}

// FieldDescription ...
func FieldDescription(value string) Field {
	return String("desc", value)
}

// FieldExtMessage ...
func FieldExtMessage(vals ...interface{}) Field {
	return zap.Any("ext", vals)
}

// FieldStack ...
func FieldStack(value []byte) Field {
	return ByteString("stack", value)
}

// FieldMethod 不推荐 未来会废弃 请使用语义更加明确的 FieldHttpMethod
func FieldMethod(value string) Field {
	return FieldHttpMethod(value)
}

// FieldHttpMethod ...
func FieldHttpMethod(value string) Field {
	return String("http.method", value)
}

// FieldDbSystem ...
func FieldDbSystem(value string) Field {
	return String("db.system", value)
}

// FieldDbOperation ...
func FieldDbOperation(value string) Field {
	return String("db.operation", value)
}

// FieldEvent ...
func FieldEvent(value string) Field {
	return String("event", value)
}

// FieldIP ...
func FieldIP(value string) Field {
	return String("ip", value)
}

// FieldCustomKeyValue 设置自定义日志
func FieldCustomKeyValue(key string, value string) Field {
	return String(strings.ToLower(key), value)
}
