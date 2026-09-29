package aliOssApi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func TestPrivateObjectStoreValidatesServerConfiguration(t *testing.T) {
	valid := PrivateObjectStoreConfig{Region: "cn-hangzhou", Endpoint: "oss-cn-hangzhou.aliyuncs.com", Bucket: "21b-erp", AccessKeyID: "test-id", AccessKeySecret: "test-secret"}
	store, err := NewPrivateObjectStore(valid)
	if err != nil || store.BucketName() != valid.Bucket {
		t.Fatalf("valid config rejected: %v", err)
	}
	for _, change := range []func(*PrivateObjectStoreConfig){func(c *PrivateObjectStoreConfig) { c.Region = "" }, func(c *PrivateObjectStoreConfig) { c.Endpoint = "http://oss-cn-hangzhou.aliyuncs.com" }, func(c *PrivateObjectStoreConfig) { c.Bucket = "bad/bucket" }, func(c *PrivateObjectStoreConfig) { c.AccessKeySecret = "" }, func(c *PrivateObjectStoreConfig) { c.Endpoint = "https://user:pass@example.com" }} {
		cfg := valid
		change(&cfg)
		if _, err := NewPrivateObjectStore(cfg); err == nil {
			t.Fatal("invalid config accepted")
		}
	}
}

func TestPrivateObjectOperationsSignV4AndKeepACLPrivate(t *testing.T) {
	var calls []string
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method)
		if !strings.HasPrefix(r.Header.Get("Authorization"), "OSS4-HMAC-SHA256 ") {
			t.Error("missing V4 signature")
		}
		if r.Header.Get("x-oss-security-token") != "test-session" {
			t.Error("session credential not forwarded")
		}
		if r.URL.Path != "/test-bucket/contracts/proof.pdf" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		switch r.Method {
		case "PUT":
			if r.Header.Get("x-oss-object-acl") != "private" || r.Header.Get("x-oss-forbid-overwrite") != "true" {
				t.Error("upload ACL or overwrite protection missing")
			}
			body, _ := io.ReadAll(r.Body)
			if string(body) != "proof" || r.Header.Get("Content-Type") != "application/pdf" {
				t.Error("upload body or MIME mismatch")
			}
			w.Header().Set("ETag", "test-etag")
			w.WriteHeader(200)
		case "GET":
			w.Header().Set("Content-Length", "5")
			fmt.Fprint(w, "proof")
		case "DELETE":
			w.WriteHeader(204)
		default:
			t.Errorf("unexpected operation %s", r.Method)
		}
	}))
	defer server.Close()
	sdkCfg := oss.LoadDefaultConfig().WithRegion("cn-hangzhou").WithEndpoint(server.URL).WithUsePathStyle(true).WithHttpClient(server.Client()).WithRetryMaxAttempts(1).WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test-id", "test-secret", "test-session"))
	store := &PrivateObjectStore{client: oss.NewClient(sdkCfg), bucket: "test-bucket"}
	ctx := context.Background()
	if err := store.Put(ctx, "contracts/proof.pdf", []byte("proof"), "application/pdf"); err != nil {
		t.Fatal(err)
	}
	body, size, err := store.Get(ctx, "contracts/proof.pdf")
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(body)
	body.Close()
	if string(data) != "proof" || size != 5 {
		t.Fatal("readback mismatch")
	}
	if err := store.Delete(ctx, "contracts/proof.pdf"); err != nil {
		t.Fatal(err)
	}
	if strings.Join(calls, ",") != "PUT,GET,DELETE" {
		t.Fatalf("unexpected operations %v", calls)
	}
	for _, key := range []string{"", "../private", "contracts/../other", "/absolute", "a\\b", "a\x00b"} {
		if err := store.Put(ctx, key, []byte("proof"), "application/pdf"); err == nil {
			t.Errorf("unsafe key accepted: %q", key)
		}
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if err := store.Delete(canceled, "contracts/proof.pdf"); err == nil {
		t.Error("canceled operation accepted")
	}
}
