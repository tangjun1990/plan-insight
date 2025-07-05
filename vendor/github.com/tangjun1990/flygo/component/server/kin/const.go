package kin

const (
	// HeaderAcceptEncoding ...
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderContentType ...
	HeaderContentType = "Content-Type"
	// HeaderGRPCPROXYError ...
	HeaderGRPCPROXYError = "GRPC-Proxy-Error"
	charsetUTF8          = "charset=utf-8"

	// MIMkapplicationJSON ...
	MIMkapplicationJSON = "application/json"
	// MIMkapplicationJSONCharsetUTF8 ...
	MIMkapplicationJSONCharsetUTF8 = MIMkapplicationJSON + "; " + charsetUTF8
	// MIMkapplicationProtobuf ...
	MIMkapplicationProtobuf = "application/protobuf"
)

const (
	codeMS                   = 1000
	codeMSInvalidParam       = 1001
	codeMSInvoke             = 1002
	codeMSInvokeLen          = 1003
	codeMSSecondItemNotError = 1004
	codeMSResErr             = 1005
)
