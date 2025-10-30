# paramHelper

`paramHelper` 把原先 `networkHelper.BuildParams` 的参数构建流程独立成一个可复用的包，负责：

- 统一从 `*http.Request` 中提取 Query、JSON、Form 数据并映射到结构体
- 通过 `validatorMng` 执行结构体验证
- 根据 `belong/kind/default` 等标签生成条件、值、附加信息等元数据
- 提供可选的校验与处理扩展点，方便在业务侧追加自定义逻辑

## 快速开始

定义参数结构体（需要实现 `networkStruct.ParamsInterface`）：

```go
type AdviceCreate struct {
    networkStruct.Params `swaggerignore:"true"`
    MerchantID uint64 `belong:"value" json:"merchant_id" validate:"required"`
    Content    string `belong:"value" json:"content" validate:"required"`
    ImgURLs    string `belong:"value" json:"img_urls" validate:"required"`
}
```

在 Handler 中调用：

```go
var params AdviceCreate
if err := paramHelper.BuildParams(r, &params, networkStruct.BodyJson); err != nil {
    networkHelper.ParamsInvalid(w, err) // 统一错误响应
    return
}

// 现在可以安全使用 params，并通过 params.GetValue()/GetCondition() 读取构建后的数据
```

> 仍需嵌入 `networkStruct.Params`，以便在 `handleParams` 阶段写回分页、排序等信息。

## 可选配置

`BuildParams` 支持可变参数 `opts ...BuildParamsOption`：

- `paramHelper.WithSkipValidation()`：跳过全部校验
- `paramHelper.WithValidators(customValidator...)`：在默认校验之后追加自定义校验
- `paramHelper.WithSkipHandle()`：跳过后续的 `handleParams`
- `paramHelper.WithMutators(customMutator...)`：在默认处理之后追加自定义处理逻辑

示例：

```go
err := paramHelper.BuildParams(r, &params, networkStruct.BodyJson,
    paramHelper.WithValidators(func(pi networkStruct.ParamsInterface) error {
        // 自定义校验
        return nil
    }),
    paramHelper.WithMutators(func(pi networkStruct.ParamsInterface) error {
        // 自定义元数据处理
        return nil
    }),
)
```

## 与 `networkHelper` 的关系

- `networkHelper` 不再内置 `BuildParams`，请直接引入 `paramHelper` 使用新接口
- 原 `fillParams/handleParams/getFormattedValue` 等实现均迁移到了 `paramHelper`
- 如需保持旧行为，请在业务侧调整依赖，避免再从 `networkHelper` 调用 `BuildParams`

## 注意事项

- 所有依赖的标签（`belong`、`kind`、`default` 等）保持与旧实现一致
- JSON 请求体会把原始字段写入 `Params.SetRawMap`，便于判断前端是否传值
- 校验失败或处理失败会返回带上下文的错误（`build params: xxx`），直接透传给调用方即可

