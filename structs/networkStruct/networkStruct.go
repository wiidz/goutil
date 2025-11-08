package networkStruct

import (
	"net/http"
)

type ContentType int8

const (
	Query    ContentType = 1
	BodyJson ContentType = 2
	BodyForm ContentType = 3
	BodyXML  ContentType = 4
	XWWWForm ContentType = 5
)

func (contentType ContentType) GetContentTypeStr() string {
	switch contentType {
	case Query:
		return ""
	case BodyJson:
		return "application/json;charset=utf-8"
	case BodyForm:
		return "application/x-www-form-urlencoded;charset=utf-8"
	case XWWWForm:
		return "application/x-www-form-urlencoded;charset=utf-8"
	case BodyXML:
		return "application/xml;charset=utf-8"
	default:
		return ""
	}
}

// ReadCommonStruct 读取列表公用的参数
type ReadCommonStruct struct {
	PageNow  int    `json:"page_now" belong:"etc" default:"1"`
	PageSize int    `json:"page_size" belong:"etc" default:"10"`
	Order    string `json:"order" belong:"etc" default:"id asc"`
}

type Method int8

const (
	Get     Method = 1
	Post    Method = 2
	Put     Method = 3
	Delete  Method = 4
	Options Method = 5
)

func (p Method) String() string {
	switch p {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	case Options:
		return "OPTIONS"
	default:
		return "UNKNOWN"
	}
}

// ParamsInterface 参数接口
type ParamsInterface interface {

	// order 顺序
	GetOrder() string
	SetOrder(string)

	// page 页码
	GetPageSize() int
	SetPageSize(int)
	GetPageNow() int
	SetPageNow(int)
	GetLimit() int
	GetOffset() int

	// [read] - prelaoad 附加
	GetPreloads() []string
	SetPreloads([]string)

	// [read、update、post、delete] - 原始数据
	GetRawMap() map[string]interface{}
	SetRawMap(map[string]interface{})

	// [read、update、delete] - condition 条件
	GetCondition() map[string]interface{}
	SetCondition(map[string]interface{})

	// [read] - rows 查询记录
	GetRows() interface{}
	SetRows(interface{})

	// [read、insert、update、delete] - row
	GetRow() interface{}
	SetRow(interface{})

	// [read] -  count 查询数目
	GetCount() int64
	SetCount(int64)

	// error 错误
	GetError() error
	SetError(error)

	// [update、delete] - rows affected 影响行数
	GetRowsAffected() int64
	SetRowsAffected(int64)

	// [update、insert] - value 数据
	GetValue() map[string]interface{}
	SetValue(map[string]interface{})

	// [update、delete] - 表名
	GetTableName() string
	SetTableName(string)

	// [insert] - 新的主键
	GetNewID() uint64
	SetNewID(newID uint64)

	// 其他
	GetEtc() map[string]interface{}
	SetEtc(data map[string]interface{})
}

// Params 参数结构体
type Params struct {

	// 前端参数
	PageNow  int    `json:"page_now" belong:"etc" default:"1"`    // [read]
	PageSize int    `json:"page_size" belong:"etc" default:"10"`  // [read]
	Order    string `json:"order" belong:"etc" default:"ids asc"` // [read]

	// 根据前端参数处理后的数据
	Condition map[string]interface{} // [read、update、delete] 条件
	Value     map[string]interface{} // [update、insert] - 数据
	Etc       map[string]interface{} // 其他
	RawMap    map[string]interface{}

	// 内部补充参数
	Single    bool     // [read] - 附加条件
	Preloads  []string // [read] - 附加
	TableName string   // [update、delete] - 指定表名

	// 操作结果
	NewID        uint64      // [insert] - 新主键
	Rows         interface{} // [read] - 结构切片
	Row          interface{} // [read、update、delete、insert] - 结构
	Count        int64       // [read] - 统计行数
	RowsAffected int64       // [update、delete] - 影响行数
	Error        error       // 错误
}

// GetTableName 表名
func (params *Params) GetTableName() string {
	return params.TableName
}
func (params *Params) SetTableName(tableName string) {
	params.TableName = tableName
}

// GetCondition 条件
func (params *Params) GetCondition() map[string]interface{} {
	return params.Condition
}
func (params *Params) SetCondition(condition map[string]interface{}) {
	params.Condition = condition
}

// GetValue Value 操作值
func (params *Params) GetValue() map[string]interface{} {
	return params.Value
}
func (params *Params) SetValue(value map[string]interface{}) {
	params.Value = value
}

// GetPageNow 页码
func (params *Params) GetPageNow() int {
	return params.PageNow
}
func (params *Params) SetPageNow(pageNow int) {
	params.PageNow = pageNow
}
func (params *Params) GetPageSize() int {
	return params.PageSize
}
func (params *Params) SetPageSize(pageSize int) {
	params.PageSize = pageSize
}
func (params *Params) GetOffset() int {
	var offset int
	if params.PageNow > 1 {
		offset = (params.PageNow - 1) * params.PageSize
	} else {
		offset = 0
	}
	return offset
}
func (params *Params) GetLimit() int {
	return params.PageSize
}

// RawMap 原始数据
func (params *Params) GetRawMap() map[string]interface{} {
	return params.RawMap
}
func (params *Params) SetRawMap(RawMap map[string]interface{}) {
	params.RawMap = RawMap
}

// Count 结果数目
func (params *Params) SetCount(count int64) {
	params.Count = count
}
func (params *Params) GetCount() int64 {
	return params.Count
}

// RowsAffected 影响行数
func (params *Params) GetRowsAffected() int64 {
	return params.RowsAffected
}
func (params *Params) SetRowsAffected(rowsAffected int64) {
	params.RowsAffected = rowsAffected
}

// Error 错误
func (params *Params) GetError() error {
	return params.Error
}
func (params *Params) SetError(err error) {
	params.Error = err
}

// Rows 查询结果
func (params *Params) SetRows(rows interface{}) {
	params.Rows = rows
}
func (params *Params) GetRows() interface{} {
	return params.Rows
}

// Row 查询结果
func (params *Params) SetRow(rows interface{}) {
	params.Row = rows
}
func (params *Params) GetRow() interface{} {
	return params.Row
}

// Order 顺序
func (params *Params) GetOrder() string {
	return params.Order
}
func (params *Params) SetOrder(order string) {
	params.Order = order
}

// Preloads 预加载
func (params *Params) GetPreloads() []string {
	return params.Preloads
}
func (params *Params) SetPreloads(preloads []string) {
	params.Preloads = preloads
}

// NewID 新主键
func (params *Params) GetNewID() uint64 {
	return params.NewID
}
func (params *Params) SetNewID(newID uint64) {
	params.NewID = newID
}

// Etc 其他
func (params *Params) GetEtc() map[string]interface{} {
	return params.Etc
}
func (params *Params) SetEtc(data map[string]interface{}) {
	params.Etc = data
}

// MyRequestParams 自己定义的请求体
// method networkStruct.Method, contentType networkStruct.ContentType, targetURL string, params map[string]interface{}, headers map[string]string, iStruct interface{}
type MyRequestParams struct {
	Method      Method // 请求方法
	URL         string // 请求地址
	ContentType ContentType
	Headers     map[string]string      // 请求头参数
	Params      map[string]interface{} // 请求参数
	ResStruct   interface{}            // 要被解析成的结构体
}

// MyRequestResp 自己定义的请求返回
// interface{}, *http.Header, int, error
type MyRequestResp struct {
	StatusCode      int
	Headers         http.Header // 返回的请求头
	ResStruct       interface{} // 解析后的结构体
	IsParsedSuccess bool        // 是否解析成功
	ResStr          string      // 返回数据字符串化
}
