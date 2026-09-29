package aliOssApi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

// PrivateObjectStoreConfig belongs on the server. It must never be serialized
// into browser responses. Region is required by OSS Signature V4.
type PrivateObjectStoreConfig struct {
	Region, Endpoint, Bucket                    string
	AccessKeyID, AccessKeySecret, SecurityToken string
}

// PrivateObjectStore provides context-aware private object operations without
// publishing permanent URLs or requiring an AssumeRole ARN. Existing OssApi
// methods keep their behavior for callers using the legacy STS workflow.
type PrivateObjectStore struct {
	client *oss.Client
	bucket string
}

func NewPrivateObjectStore(cfg PrivateObjectStoreConfig) (*PrivateObjectStore, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if !strings.Contains(endpoint, "://") {
		endpoint = "https://" + endpoint
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("OSS endpoint must be an HTTPS service endpoint")
	}
	if !regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,61}[a-z0-9]$`).MatchString(cfg.Bucket) {
		return nil, errors.New("invalid OSS bucket name")
	}
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(cfg.Region) {
		return nil, errors.New("OSS region is required")
	}
	if strings.TrimSpace(cfg.AccessKeyID) == "" || strings.TrimSpace(cfg.AccessKeySecret) == "" {
		return nil, errors.New("OSS credentials are required")
	}
	sdkCfg := oss.LoadDefaultConfig().WithRegion(cfg.Region).WithEndpoint(endpoint).
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, cfg.SecurityToken)).
		WithRetryMaxAttempts(2).WithConnectTimeout(10 * time.Second).WithReadWriteTimeout(60 * time.Second)
	return &PrivateObjectStore{client: oss.NewClient(sdkCfg), bucket: cfg.Bucket}, nil
}

func (s *PrivateObjectStore) BucketName() string { return s.bucket }

// Put uploads bytes with a private object ACL and refuses to overwrite an
// existing key. The caller chooses a randomized, authorized object key.
func (s *PrivateObjectStore) Put(ctx context.Context, key string, content []byte, contentType string) error {
	if err := validPrivateObjectKey(key); err != nil {
		return err
	}
	_, err := s.client.PutObject(ctx, &oss.PutObjectRequest{Bucket: oss.Ptr(s.bucket), Key: oss.Ptr(key), Body: bytes.NewReader(content), ContentType: oss.Ptr(contentType), ContentLength: oss.Ptr(int64(len(content))), Acl: oss.ObjectACLPrivate, ForbidOverwrite: oss.Ptr("true")})
	return err
}

// Get streams an authenticated object. The caller must close the returned body.
func (s *PrivateObjectStore) Get(ctx context.Context, key string) (io.ReadCloser, int64, error) {
	if err := validPrivateObjectKey(key); err != nil {
		return nil, 0, err
	}
	result, err := s.client.GetObject(ctx, &oss.GetObjectRequest{Bucket: oss.Ptr(s.bucket), Key: oss.Ptr(key)})
	if err != nil {
		return nil, 0, err
	}
	return result.Body, result.ContentLength, nil
}

func (s *PrivateObjectStore) Delete(ctx context.Context, key string) error {
	if err := validPrivateObjectKey(key); err != nil {
		return err
	}
	_, err := s.client.DeleteObject(ctx, &oss.DeleteObjectRequest{Bucket: oss.Ptr(s.bucket), Key: oss.Ptr(key)})
	return err
}

func validPrivateObjectKey(key string) error {
	if key == "" || len(key) > 1023 || strings.HasPrefix(key, "/") || strings.ContainsAny(key, "\\\x00\r\n") {
		return errors.New("invalid OSS object key")
	}
	for _, part := range strings.Split(key, "/") {
		if part == "" || part == "." || part == ".." {
			return errors.New("invalid OSS object key")
		}
	}
	return nil
}

// PrivateObjectErrorCode provides diagnostic codes without exposing credentials,
// signed request URLs, or sensitive content from upstream error messages.
func PrivateObjectErrorCode(err error) string {
	var serviceError *oss.ServiceError
	if errors.As(err, &serviceError) {
		return fmt.Sprintf("status=%d code=%s", serviceError.StatusCode, serviceError.Code)
	}
	return "connection or configuration error"
}
