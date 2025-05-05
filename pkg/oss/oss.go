package oss

import (
	"context"
	"errors"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"strings"
)

var (
	ErrInvalidURLPrefix = errors.New("文件地址路径异常 应以http/https开头")
	ErrInvalidURLPath   = errors.New("文件地址异常 请正确上传文件")
	ErrGetObject        = errors.New("文件获取异常 请稍后重试或联系客服处理")
)

type Bucket struct {
	accessKeyId     string
	accessKeySecret string
	endpoint        string

	bucketNameFunc func() string
	pathFunc       func(fileName string) string

	client *oss.Bucket
}

// NewBucket .
func NewBucket(accessKeyId, accessKeySecret, endpoint string, bucketNameFun func() string, pathFun func(fileName string) string) *Bucket {
	return &Bucket{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		endpoint:        endpoint,
		bucketNameFunc:  bucketNameFun,
		pathFunc:        pathFun,
	}
}

func (b *Bucket) initBucket() error {
	if b.client != nil {
		return nil
	}

	client, err := oss.New("https://"+b.endpoint, b.accessKeyId, b.accessKeySecret)
	if err != nil {
		return err
	}
	c, err := client.Bucket(b.bucketNameFunc())
	if err != nil {
		return err
	}
	b.client = c
	return nil
}

// IsExisted 查询文件是否存在
func (b *Bucket) IsExisted(fileName string) (bool, error) {
	err := b.initBucket()
	if err != nil {
		return false, err
	}

	return b.client.IsObjectExist(b.pathFunc(fileName))
}

// Upload 上传文件到 OSS
// return 文件访问地址
func (b *Bucket) Upload(fileName string, reader io.Reader) (string, error) {
	err := b.initBucket()
	if err != nil {
		return "", err
	}
	err = b.client.PutObject(b.pathFunc(fileName), reader, oss.CacheControl("public"))
	if err != nil {
		return "", err
	}

	return b.FullPath(fileName), nil
}

// FullPath oss全路径
func (b *Bucket) FullPath(fileName string) string {
	return fmt.Sprintf("https://%s.%s/%s", b.bucketNameFunc(), b.endpoint, b.pathFunc(fileName))
}

// UrlToObjectKey 获取 objectKey, FullPath 的逆方法
func (b *Bucket) UrlToObjectKey(url string) (string, error) {
	httpStr := "http"
	if !strings.HasPrefix(url, httpStr) {
		return "", ErrInvalidURLPrefix
	}
	if strings.HasPrefix(url, "https") {
		httpStr = "https"
	}

	prefix := fmt.Sprintf("%s://%s.%s/", httpStr, b.bucketNameFunc(), b.endpoint)
	key := strings.TrimPrefix(url, prefix)
	if strings.Contains(key, "://") {
		return "", ErrInvalidURLPath
	}

	return key, nil
}

// GetFileWithURL 下载文件
func (b *Bucket) GetFileWithURL(ctx context.Context, url string, options ...oss.Option) (*oss.GetObjectResult, error) {
	objectKey, err := b.UrlToObjectKey(url)
	if err != nil {
		return nil, err
	}
	err = b.initBucket()
	if err != nil {
		return nil, err
	}

	res, err := b.client.DoGetObject(&oss.GetObjectRequest{ObjectKey: objectKey}, options)
	if err != nil {
		return nil, ErrGetObject
	}
	return res, nil
}
