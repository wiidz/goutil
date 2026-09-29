# 私有 OSS 文件存储

`NewPrivateObjectStore` 为服务器端提供私有文件上传、流式读取和删除能力，使用阿里云 Go SDK V2 与 Signature V4。原有 `NewOssApi`、STS、上传签名和文件路径上传方法继续保留。

```go
store, err := aliOssApi.NewPrivateObjectStore(aliOssApi.PrivateObjectStoreConfig{
    Region:          "cn-hangzhou",
    Endpoint:        "https://oss-cn-hangzhou.aliyuncs.com",
    Bucket:          "21b-erp",
    AccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
    AccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
})
if err != nil {
    return err
}
if err := store.Put(ctx, "contracts/random-id.pdf", content, "application/pdf"); err != nil {
    return err
}
body, size, err := store.Get(ctx, "contracts/random-id.pdf")
if err != nil {
    return err
}
// 读取 body / size 后，关闭流，再按业务需要删除对象。
_ = size
body.Close()
err = store.Delete(ctx, "contracts/random-id.pdf")
```

配置要求 HTTPS Endpoint、地域、Bucket 和访问凭证；已有临时凭证可通过 `SecurityToken` 传入，不要求额外配置 AssumeRole ARN。临时凭证的刷新由调用方管理。凭证只保留在服务器端，不得通过业务接口序列化给浏览器。

上传明确设置对象 ACL 为 `private`，并禁止覆盖同名对象。调用方负责生成随机文件名、验证文件内容和大小、校验业务权限，并在数据库记录失败时清理已上传对象。对象读取由应用服务器完成，便于每次下载都检查业务权限。上传、读取、删除都支持 `context` 取消。

诊断时可以调用 `PrivateObjectErrorCode(err)`，它只返回 OSS 错误码和状态码；不要把原始 SDK 错误直接显示给用户或记录到日志，避免输出已签名 URL 等敏感信息。

验证命令：

```sh
go test ./apis/ali/aliOssApi
go vet ./apis/ali/aliOssApi
```

本次 OSS 能力已合入 `v1.5.99` 对应的最新 `master` 基线，保留该版本的 Redis 包路径、认证支持和 `database` 配置字段。CRM 当前正式依赖仍为 `goutil v1.5.99`，开发联调使用该版本加上本次新增方法。发布包含这次提交的新版本后，CRM API 需要升级到该实际版本并验证不依赖开发 workspace 的构建。
