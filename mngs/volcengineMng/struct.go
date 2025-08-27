package volcengineMng

import (
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/wiidz/goutil/structs/configStruct"
	"strings"
)

type Config struct {
	ApiKey string
}

type VolcengineMng struct {
	Config *configStruct.VolcengineConfig
	Client *arkruntime.Client
}

type AIModel string

const (
	Doubao AIModel = "doubao-1-5-pro-32k-250115"
)

// IsVisionModel 检查是否是视觉模型
func (m AIModel) IsVisionModel() bool {
	return strings.Contains(string(m), "vision")
}

type Role string

const (
	User   Role = "user"
	System Role = "system"
)

type ChatParam struct {
	Role Role
	Text string
}

type ThinkingType string

const (
	Disabled ThinkingType = "disabled"
	Enabled  ThinkingType = "enabled"
	Auto     ThinkingType = "auto"
)

func (m ThinkingType) GetThinkingType() model.ThinkingType {
	if m == Disabled {
		return model.ThinkingTypeDisabled
	} else if m == Enabled {
		return model.ThinkingTypeEnabled
	} else if m == Auto {
		return model.ThinkingTypeAuto
	} else {
		return model.ThinkingTypeDisabled
	}
}

// SearchWikiResp 搜索知识库相应
type SearchWikiResp struct {
	Code int `json:"code"`
	Data struct {
		CollectionName string `json:"collection_name"`
		Count          int    `json:"count"`
		RewriteQuery   string `json:"rewrite_query"`
		TokenUsage     struct {
			EmbeddingTokenUsage struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			} `json:"embedding_token_usage"`
			RerankTokenUsage int `json:"rerank_token_usage"`
		} `json:"token_usage"`
		ResultList []struct {
			Id          string  `json:"id"`
			Content     string  `json:"content"`
			Score       float64 `json:"score"`
			PointId     string  `json:"point_id"`
			ChunkTitle  string  `json:"chunk_title"`
			ChunkId     int     `json:"chunk_id"`
			ProcessTime int     `json:"process_time"`
			DocInfo     struct {
				DocId      string `json:"doc_id"`
				DocName    string `json:"doc_name"`
				CreateTime int    `json:"create_time"`
				DocType    string `json:"doc_type"`
				DocMeta    string `json:"doc_meta"`
				Source     string `json:"source"`
				Title      string `json:"title"`
			} `json:"doc_info"`
			RecallPosition int    `json:"recall_position"`
			ChunkType      string `json:"chunk_type"`
		} `json:"result_list"`
	} `json:"data"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
}

type SearchWikiParam struct {
	Name           string          `json:"name"`
	Query          string          `json:"query"`
	Limit          int             `json:"limit"`
	QueryParam     interface{}     `json:"query_param"` // Assuming query_param is a map or similar
	DenseWeight    float64         `json:"dense_weight"`
	PreProcessing  *PreProcessing  `json:"pre_processing"`
	PostProcessing *PostProcessing `json:"post_processing"`
}

type PostProcessing struct {
	RerankSwitch      bool   `json:"rerank_switch"`
	RerankModel       string `json:"rerank_model"`
	RerankOnlyChunk   bool   `json:"rerank_only_chunk"`
	RetrieveCount     int    `json:"retrieve_count"`
	EndpointID        string `json:"endpoint_id"`
	ChunkGroup        bool   `json:"chunk_group"`
	GetAttachmentLink bool   `json:"get_attachment_link"`
}

type PreProcessing struct {
	NeedInstruction  bool       `json:"need_instruction"`
	Rewrite          bool       `json:"rewrite"`
	Messages         []*Message `json:"messages"`
	ReturnTokenUsage bool       `json:"return_token_usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
