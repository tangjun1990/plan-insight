package oss

import (
	"fmt"

	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/plan-insight/pkg/util"
)

const (
	tempAccessKeyId     = "LTAI5t8hxys6L8jvjBQ5USFC"
	tempAccessKeySecret = "9LffuI8i40ilii6wGp0c1lFLs31H88"
	tempEndpoint        = "oss-cn-shanghai.aliyuncs.com"
)

func NewTempBucket() *Bucket {
	return NewBucket(
		tempAccessKeyId,
		tempAccessKeySecret,
		tempEndpoint,
		tempBucketName,
		tempPath)
}

func tempPath(fileName string) string {
	return fmt.Sprintf("temp/%s/%s", util.NowMonth(), fileName)
}

func tempBucketName() string {
	if kcfg.GetString("app.global.environment") == "prod" {
		return "febooksh"
	}
	return "febookshtest"
}
