package aiMng

import "context"

// ProductPriceRecord 产品价格记录（通用类型，用于 AI 管理器）
type ProductPriceRecord struct {
	ID          uint64            `json:"id"`
	SpuID       uint64            `json:"spu_id,omitempty"`
	SKUCode     string            `json:"sku_code"`
	DisplayName string            `json:"display_name"`
	PriceIn     *float64          `json:"price_in,omitempty"`
	PricePromot *float64          `json:"price_promot,omitempty"`
	Currency    string            `json:"currency,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	Thumbnail   string            `json:"thumbnail_img,omitempty"`
}

// ProductPriceQueryParams 产品价格查询参数
type ProductPriceQueryParams struct {
	ID         uint64
	SearchText string
	Limit      int
}

// ProductQueryDebug 产品查询调试信息
type ProductQueryDebug struct {
	SQL        string                   `json:"sql"`
	SearchText string                   `json:"search_text"`
	Normalized string                   `json:"normalized"`
	Segments   [][]string               `json:"segments"`
	Samples    []map[string]interface{} `json:"samples,omitempty"`
}

// ProductPriceQuery 产品价格查询接口（由业务代码实现）
type ProductPriceQuery interface {
	QueryProductPrices(ctx context.Context, params ProductPriceQueryParams) ([]ProductPriceRecord, int64, *ProductQueryDebug, error)
}

