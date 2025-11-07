// Package wikiMng 提供火山引擎知识库管理功能
// 包含知识库检索、对话生成等核心功能的类型定义和工具函数
package wikiMng

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/wiidz/goutil/mngs/volcengineMng"
)

// 系统字段常量定义
const (
	// SysFieldDocName 文档名称字段
	SysFieldDocName = "doc_name"
	// SysFieldTitle 文档标题字段
	SysFieldTitle = "title"
	// SysFieldChunkTitle 文档切片标题字段
	SysFieldChunkTitle = "chunk_title"
	// SysFieldContent 文档内容字段（必传）
	SysFieldContent = "content"
)

// WikiMng 火山引擎知识库管理器
type WikiMng struct {
	Config *Config
}

// Config 配置信息
type Config struct {
	AccessKeyID string // Access Key
	SecretKey   string // Secret Key
	ApiKey      string // API Key

	StreamTimeout time.Duration // 流式超时
	SimpleTimeout time.Duration // 单此请求
}

type AIModelSetting struct {
	Model        volcengineMng.AIModel `json:"model"`         // 模型名称
	ModelVersion string                `json:"model_version"` // 模型版本
	MaxTokens    int32                 `json:"max_tokens"`    // 最大Token数
	Temperature  float32               `json:"temperature"`   // 温度参数
}

// PromptExtraContext 提示词额外上下文配置
// 用于自定义拼接Prompt时传入的字段
//
// 使用建议：
//  1. 非结构化数据：SystemFields 中可传入 doc_name, title, chunk_title, content 等字段
//     其中 content 字段为必传字段，其他为可选字段
//  2. 结构化数据：SystemFields 中可传入 title 字段(可选)，
//     SelfDefineFields 字段来源于结构化数据的表头字段，
//     表头字段中索引字段需要必传，非索引字段可以不传
type PromptExtraContext struct {
	SelfDefineFields []string `json:"self_define_fields"` // 自定义字段列表
	SystemFields     []string `json:"system_fields"`      // 系统字段列表
}

// CollectionSearchKnowledgeRequest 知识库检索请求参数
type CollectionSearchKnowledgeRequest struct {
	Name           string          `json:"name,omitempty"`            // 知识库名称
	Project        string          `json:"project,omitempty"`         // 项目名称
	ResourceId     string          `json:"resource_id,omitempty"`     // 资源ID
	Query          string          `json:"query"`                     // 查询内容
	Limit          int32           `json:"limit"`                     // 返回结果数量限制
	QueryParam     *QueryParamInfo `json:"query_param"`               // 查询参数
	DenseWeight    float32         `json:"dense_weight"`              // 密集检索权重
	MdSearch       bool            `json:"md_search"`                 // 是否启用Markdown搜索
	Preprocessing  *PreProcessing  `json:"pre_processing,omitempty"`  // 预处理参数
	Postprocessing *PostProcessing `json:"post_processing,omitempty"` // 后处理参数
}

// QueryParamInfo 查询参数信息
type QueryParamInfo struct {
	DocFilter interface{} `json:"doc_filter"` // 文档过滤器
}

// MessageParam 消息参数
type MessageParam struct {
	Role    string      `json:"role"`    // 角色：system/user/assistant
	Content interface{} `json:"content"` // 消息内容
}

// ChatCompletionMessageContent 聊天完成消息内容
type ChatCompletionMessageContent struct {
	StringValue *string                             // 字符串值
	ListValue   []*ChatCompletionMessageContentPart // 列表值（多模态）
}

// ChatMessageImageURL 聊天消息图片URL
type ChatMessageImageURL struct {
	URL string `json:"url,omitempty"` // 图片URL
}

// ChatCompletionMessageContentPartType 消息内容部分类型
type ChatCompletionMessageContentPartType string

const (
	// ChatCompletionMessageContentPartTypeText 文本类型
	ChatCompletionMessageContentPartTypeText ChatCompletionMessageContentPartType = "text"
	// ChatCompletionMessageContentPartTypeImageURL 图片URL类型
	ChatCompletionMessageContentPartTypeImageURL ChatCompletionMessageContentPartType = "image_url"
)

// ChatCompletionMessageContentPart 聊天完成消息内容部分
type ChatCompletionMessageContentPart struct {
	Type     ChatCompletionMessageContentPartType `json:"type,omitempty"`      // 内容类型
	Text     string                               `json:"text,omitempty"`      // 文本内容
	ImageURL *ChatMessageImageURL                 `json:"image_url,omitempty"` // 图片URL
}

// PreProcessing 检索接口预处理参数
type PreProcessing struct {
	NeedInstruction  bool            `json:"need_instruction"`   // 是否需要指令
	Rewrite          bool            `json:"rewrite"`            // 是否重写查询
	Messages         []*MessageParam `json:"messages"`           // 消息列表
	ReturnTokenUsage bool            `json:"return_token_usage"` // 是否返回token使用情况
}

// PostProcessing 检索接口后处理参数
type PostProcessing struct {
	RerankSwitch        bool                   `json:"rerank_switch"`                   // 重排序开关
	RerankModel         string                 `json:"rerank_model,omitempty"`          // 重排序模型
	RerankOnlyChunk     bool                   `json:"rerank_only_chunk"`               // 仅重排序切片
	RetrieveCount       int32                  `json:"retrieve_count"`                  // 检索数量
	EndpointID          string                 `json:"endpoint_id"`                     // 端点ID
	ChunkDiffusionCount int32                  `json:"chunk_diffusion_count"`           // 切片扩散数量
	ChunkGroup          bool                   `json:"chunk_group"`                     // 切片分组
	ChunkScoreAggrType  string                 `json:"chunk_score_aggr_type,omitempty"` // 切片分数聚合类型
	ChunkExtraContent   map[string]interface{} `json:"chunk_extra_content"`             // 切片额外内容
	GetAttachmentLink   bool                   `json:"get_attachment_link"`             // 是否获取附件链接
}

// CollectionSearchKnowledgeResponse 知识库检索响应
type CollectionSearchKnowledgeResponse struct {
	Code    int64                                  `json:"code"`    // 响应码
	Message string                                 `json:"message"` // 响应消息
	Data    *CollectionSearchKnowledgeResponseData `json:"data"`    // 响应数据
}

// CollectionSearchKnowledgeResponseData 检索响应数据
type CollectionSearchKnowledgeResponseData struct {
	CollectionName string                          `json:"collection_name"`         // 集合名称
	Count          int32                           `json:"count"`                   // 结果数量
	RewriteQuery   string                          `json:"rewrite_query,omitempty"` // 重写后的查询
	TokenUsage     *TotalTokenUsage                `json:"token_usage,omitempty"`   // Token使用情况
	ResultList     []*CollectionSearchResponseItem `json:"result_list,omitempty"`   // 结果列表
}

// TotalTokenUsage 总Token使用情况
type TotalTokenUsage struct {
	EmbeddingUsage *ModelTokenUsage `json:"embedding_token_usage,omitempty"` // 嵌入模型Token使用
	RerankUsage    *int64           `json:"rerank_token_usage,omitempty"`    // 重排序Token使用
	LLMUsage       *ModelTokenUsage `json:"llm_token_usage,omitempty"`       // 大语言模型Token使用
	RewriteUsage   *ModelTokenUsage `json:"rewrite_token_usage,omitempty"`   // 重写Token使用
}

// CollectionSearchResponseItem 检索响应项
type CollectionSearchResponseItem struct {
	Id                  string                              `json:"id"`                            // 项目ID
	Content             string                              `json:"content"`                       // 内容
	MdContent           string                              `json:"md_content,omitempty"`          // Markdown内容
	Score               float64                             `json:"score"`                         // 分数
	PointId             string                              `json:"point_id"`                      // 点ID
	OriginText          string                              `json:"origin_text,omitempty"`         // 原始文本
	OriginalQuestion    string                              `json:"original_question,omitempty"`   // 原始问题
	ChunkTitle          string                              `json:"chunk_title,omitempty"`         // 切片标题
	ChunkId             int                                 `json:"chunk_id"`                      // 切片ID
	ProcessTime         int64                               `json:"process_time"`                  // 处理时间
	RerankScore         float64                             `json:"rerank_score,omitempty"`        // 重排序分数
	DocInfo             CollectionSearchResponseItemDocInfo `json:"doc_info,omitempty"`            // 文档信息
	RecallPosition      int32                               `json:"recall_position"`               // 召回位置
	RerankPosition      int32                               `json:"rerank_position,omitempty"`     // 重排序位置
	ChunkType           string                              `json:"chunk_type,omitempty"`          // 切片类型
	ChunkSource         string                              `json:"chunk_source,omitempty"`        // 切片来源
	UpdateTime          int64                               `json:"update_time"`                   // 更新时间
	ChunkAttachmentList []ChunkAttachment                   `json:"chunk_attachment,omitempty"`    // 切片附件列表
	TableChunkFields    []PointTableChunkField              `json:"table_chunk_fields,omitempty"`  // 表格切片字段
	OriginalCoordinate  *ChunkPositions                     `json:"original_coordinate,omitempty"` // 原始坐标
}

// CollectionSearchResponseItemDocInfo 检索响应项文档信息
type CollectionSearchResponseItemDocInfo struct {
	Docid      string `json:"doc_id"`             // 文档ID
	DocName    string `json:"doc_name"`           // 文档名称
	CreateTime int64  `json:"create_time"`        // 创建时间
	DocType    string `json:"doc_type"`           // 文档类型
	DocMeta    string `json:"doc_meta,omitempty"` // 文档元数据
	Source     string `json:"source"`             // 来源
	Title      string `json:"title,omitempty"`    // 标题
}

// ChunkAttachment 切片附件
type ChunkAttachment struct {
	UUID    string `json:"uuid,omitempty"` // 唯一标识
	Caption string `json:"caption"`        // 标题
	Type    string `json:"type"`           // 类型
	Link    string `json:"link,omitempty"` // 链接
}

// PointTableChunkField 点表格切片字段
type PointTableChunkField struct {
	FieldName  string      `json:"field_name"`  // 字段名称
	FieldValue interface{} `json:"field_value"` // 字段值
}

// ChunkPositions 切片位置信息
type ChunkPositions struct {
	PageNo []int       `json:"page_no"` // 页码列表
	BBox   [][]float64 `json:"bbox"`    // 边界框坐标
}

// CollectionChatCompletionRequest 聊天完成请求
type CollectionChatCompletionRequest struct {
	Model            string          `json:"model"`              // 模型名称
	ModelVersion     string          `json:"model_version"`      // 模型版本
	APIKey           string          `json:"api_key"`            // API密钥
	MaxTokens        int32           `json:"max_tokens"`         // 最大Token数
	Temperature      float32         `json:"temperature"`        // 温度参数
	Messages         []*MessageParam `json:"messages"`           // 消息列表
	Stream           bool            `json:"stream"`             // 是否流式
	ReturnTokenUsage bool            `json:"return_token_usage"` // 是否返回Token使用情况;``
}

// CollectionChatCompletionResponse 聊天完成响应
type CollectionChatCompletionResponse struct {
	Code    int64                                 `json:"code"`    // 响应码
	Message string                                `json:"message"` // 响应消息
	Data    *CollectionChatCompletionResponseData `json:"data"`    // 响应数据
}

// CollectionChatCompletionResponseData 聊天完成响应数据
type CollectionChatCompletionResponseData struct {
	GenerateAnswer   string `json:"generated_answer"`            // 生成的答案
	Usage            string `json:"usage"`                       // 使用情况
	ReasoningContent string `json:"reasoning_content,omitempty"` // 推理内容（仅推理模型）
}

// ModelTokenUsage 模型Token使用情况
type ModelTokenUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`     // 提示Token数
	CompletionTokens int64 `json:"completion_tokens"` // 完成Token数
	TotalTokens      int64 `json:"total_tokens"`      // 总Token数
}

// ParseJsonUseNumber 使用Number类型解析JSON
// 避免大数字精度丢失问题
func ParseJsonUseNumber(input []byte, target interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(input))
	if decoder == nil {
		return fmt.Errorf("failed to create JSON decoder")
	}
	decoder.UseNumber()
	if err := decoder.Decode(&target); err != nil {
		return fmt.Errorf("JSON decode failed: %w", err)
	}
	return nil
}

// SerializeToJsonBytesUseNumber 将对象序列化为JSON字节数组
// 使用Number类型避免大数字精度丢失
func (mng *WikiMng) serializeToJsonBytesUseNumber(source interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(source); err != nil {
		return nil, fmt.Errorf("JSON encode failed: %w", err)
	}
	return buf.Bytes(), nil
}

// scanDoubleCRLF 自定义分隔符函数，用于分隔 \r\n\r\n
// 主要用于流式响应的解析
func scanDoubleCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 查找 \r\n\r\n 分隔符
	if i := bytes.Index(data, []byte("\r\n\r\n")); i >= 0 {
		return i + 4, data[0:i], nil
	}
	// 处理文件结束且包含结束标记的情况
	if atEOF && strings.Contains(string(data), "\"end\":true") {
		return len(data), data, nil
	}
	return 0, nil, nil
}
