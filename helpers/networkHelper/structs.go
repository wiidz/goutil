package networkHelper

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

	GetRawJsonMap()map[string]interface{}
	SetRawJsonMap(map[string]interface{})

	// [read、update、delete] - condition 条件
	GetCondition() map[string]interface{}
	SetCondition(map[string]interface{})

	// [read] - rows 查询记录
	GetRows() interface{}
	SetRows(interface{})

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
	GetValue()map[string]interface{}
	SetValue(map[string]interface{})

	// [update、delete] - 表名
	GetTableName()string
	SetTableName(string)

	// [insert] - 新的主键
	GetNewID(ID uint64)
}

// Params 参数结构体
type Params struct {

	// 前端参数
	PageNow   int    `json:"page_now" belong:"etc" default:"1"` // [read]
	PageSize  int    `json:"page_size" belong:"etc" default:"10"` // [read]
	Order     string `json:"order" belong:"etc" default:"ids asc"` // [read]

	// 根据前端参数处理后的数据
	Condition map[string]interface{} // [read、update、delete] 条件
	Value      map[string]interface{} // [update、insert] - 数据
	RawJsonMap map[string]interface{}

	// 内部补充参数
	Single    bool // [read] - 附加条件
	Preloads  []string // [read] - 附加
	TableName  string // [update、delete] - 指定表名

	// 操作结果
	NewID uint64 // [insert] - 新主键
	Rows      interface{} // [read] - 查询结构
	Count     int64 // [read] - 统计行数
	RowsAffected int64 // [update、delete] - 影响行数
	Error        error // 错误
}

// TableName 表名
func (params *Params) GetTableName() string {
	return params.TableName
}
func (params *Params) SetTableName(tableName string) {
	params.TableName = tableName
}

// Condition 条件
func (params *Params) GetCondition() map[string]interface{} {
	return params.Condition
}
func (params *Params) SetCondition(condition map[string]interface{}) {
	params.Condition = condition
}

// Value 操作值
func (params *Params) GetValue() map[string]interface{} {
	return params.Value
}
func (params *Params) SetValue(value map[string]interface{}) {
	params.Value = value
}

// Page 页码
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

// RawJsonMap 原始数据
func (params *Params) GetRawJsonMap() map[string]interface{} {
	return params.RawJsonMap
}
func (params *Params) SetRawJsonMap(rawJsonMap map[string]interface{}) {
	params.RawJsonMap = rawJsonMap
}


// Count 结果数目
func (params *Params) SetCount(count int64) {
	params.Count = count
}
func (params *Params) GetCount() int64{
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
func (params *Params) GetRows(rows interface{}) {
	params.Rows = rows
}
func (params *Params) SetRows() interface{} {
	return params.Rows
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
func (params *Params) SetPreloads(preloads []string)  {
  params.Preloads = preloads
}

// NewID 新主键
func (params *Params) GetNewID() uint64 {
	return params.NewID
}
func (params *Params) SetNewID(newID uint64) {
	params.NewID = newID
}