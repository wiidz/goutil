package mysqlMng

// NullType 用于判断是否是null值
type NullType byte

const (
	_         NullType = iota
	IsNull             // IsNull the same as `is null`
	IsNotNull          // IsNotNull the same as `is not null`
)

// OnlyID 简单方法种用来获取新插入记录的ID值
type OnlyID struct {
	ID uint64 `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null" json:"id"` // 编号
}

// IDInterface 结构化方法种用来获取新插入记录的ID值
type IDInterface interface {
	GetID() uint64
}

type BaseInterface interface {
	GetError() error
	GetRowsAffected() int64
	SetError(error)
	SetRowsAffected(int64)
}

type BaseStruct struct {
	Error        error
	RowsAffected int64
}

// ReadInterface 读接口
type ReadInterface interface {
	BaseInterface
	// Getter
	GetPreloads() []string
	GetCondition() map[string]interface{}
	GetOrder() string
	GetOffset() int
	GetPageSize() int
	GetPageNow() int
	GetRow() interface{}
	GetRows() interface{}
	GetCount() int64
	GetLimit() int
	//Setter
	SetRow(interface{})
	SetRows(interface{})
	SetCount(int64)
	SetCondition(map[string]interface{})
	SetOrder(string)
	SetPageNow(int)
	SetPageSize(int)
}

// Read 读取操作的基本结构体
type Read struct {
	BaseStruct
	Condition map[string]interface{}
	PageNow   int    `json:"page_now" belong:"etc" default:"1"`
	PageSize  int    `json:"page_size" belong:"etc" default:"10"`
	Order     string `json:"order" belong:"etc" default:"ids asc"`
	Single    bool
	Preloads  []string
	Rows      interface{}
	Count     int64
}

func (read *Read) GetOrder() string {
	return read.Order
}
func (read *Read) GetRows() interface{} {
	return read.Rows
}
func (read *Read) GetCount() int64 {
	return read.Count
}
func (read *Read) GetCondition() map[string]interface{} {
	return read.Condition
}
func (read *Read) GetOffset() int {
	var offset int
	if read.PageNow > 1 {
		offset = (read.PageNow - 1) * read.PageSize
	} else {
		offset = 0
	}
	return offset
}
func (read *Read) GetLimit() int {
	return read.PageSize
}
func (read *Read) GetPageSize() int {
	return read.PageSize
}
func (read *Read) GetPageNow() int {
	return read.PageNow
}
func (read *Read) GetPreloads() []string {
	return read.Preloads
}
func (read *Read) GetRowsAffected() int64 {
	return read.RowsAffected
}
func (read *Read) GetError() error {
	return read.Error
}

func (read *Read) SetCount(count int64) {
	read.Count = count
}
func (read *Read) SetCondition(condition map[string]interface{}) {
	read.Condition = condition
}
func (read *Read) SetOrder(order string) {
	read.Order = order
}
func (read *Read) SetRows(rows interface{}) {
	read.Rows = rows
}
func (read *Read) SetPageNow(pageNow int) {
	read.PageNow = pageNow
}
func (read *Read) SetPageSize(pageSize int) {
	read.PageSize = pageSize
}
func (read *Read) SetError(err error) {
	read.Error = err
}
func (read *Read) SetRowsAffected(rowsAffected int64) {
	read.RowsAffected = rowsAffected
}

// JsonInterface 涉及到json操作的接口
type JsonInterface interface {
	// Getter
	GetRawMap() map[string]interface{}
	//Setter
	SetRawMap(map[string]interface{})
}

// InsertInterface 插入接口
type InsertInterface interface {
	JsonInterface
	BaseInterface
	// Getter
	GetRow() interface{}
	GetNewID() uint64
	//Setter
	SetRow(interface{})
	SetNewID(uint64)
}

// Insert 插入操作的基本结构体
type Insert struct {
	BaseStruct
	NewID  uint64 // 没有实际意义，无法实现，只为了区分interface
	Row    interface{}
	RawMap map[string]interface{}
}

func (insert *Insert) GetRawMap() map[string]interface{} {
	return insert.RawMap
}
func (insert *Insert) GetRow() interface{} {
	return insert.Row
}
func (insert *Insert) GetNewID() uint64 {
	return insert.NewID
}
func (insert *Insert) GetRowsAffected() int64 {
	return insert.RowsAffected
}
func (insert *Insert) GetError() error {
	return insert.Error
}
func (insert *Insert) SetRawMap(rawJsonMap map[string]interface{}) {
	insert.RawMap = rawJsonMap
}
func (insert *Insert) SetRow(row interface{}) {
	insert.Row = row
}
func (insert *Insert) SetRowsAffected(rowsAffected int64) {
	insert.RowsAffected = rowsAffected
}
func (insert *Insert) SetNewID(newID uint64) {
	insert.NewID = newID
}
func (insert *Insert) SetError(err error) {
	insert.Error = err
}

// UpdateInterface 修改接口
type UpdateInterface interface {
	JsonInterface
	BaseInterface
	// Getter
	GetTableName() string
	GetCondition() map[string]interface{}
	GetValue() map[string]interface{}
	GetRow() interface{}
	//Setter
	SetTableName(string)
	SetCondition(map[string]interface{})
	SetValue(map[string]interface{})
	SetRow(interface{})
}

// Update 修改操作的基本结构体
type Update struct {
	BaseStruct
	TableName string
	Condition map[string]interface{}
	Value     map[string]interface{}
	RawMap    map[string]interface{}
}

func (update *Update) GetRawMap() map[string]interface{} {
	return update.RawMap
}
func (update *Update) GetTableName() string {
	return update.TableName
}
func (update *Update) GetCondition() map[string]interface{} {
	return update.Condition
}
func (update *Update) GetValue() map[string]interface{} {
	return update.Value
}
func (update *Update) GetRowsAffected() int64 {
	return update.RowsAffected
}
func (update *Update) GetError() error {
	return update.Error
}
func (update *Update) SetRawMap(rawJsonMap map[string]interface{}) {
	update.RawMap = rawJsonMap
}
func (update *Update) SetTableName(tableName string) {
	update.TableName = tableName
}
func (update *Update) SetCondition(condition map[string]interface{}) {
	update.Condition = condition
}
func (update *Update) SetValue(value map[string]interface{}) {
	update.Value = value
}
func (update *Update) SetRowsAffected(rowsAffected int64) {
	update.RowsAffected = rowsAffected
}
func (update *Update) SetError(err error) {
	update.Error = err
}

// DeleteInterface 删除接口
type DeleteInterface interface {
	JsonInterface
	BaseInterface
	// Getter
	GetCondition() map[string]interface{}
	GetRow() interface{}
	//Setter
	SetCondition(map[string]interface{})
	SetRow(interface{})
}

// Delete 删除操作的基本结构体
type Delete struct {
	BaseStruct
	Condition map[string]interface{}
	Row       interface{}
	RawMap    map[string]interface{}
}

func (delete *Delete) GetRawMap() map[string]interface{} {
	return delete.RawMap
}
func (delete *Delete) GetCondition() map[string]interface{} {
	return delete.Condition
}
func (delete *Delete) GetRow() interface{} {
	return delete.Row
}
func (delete *Delete) GetRowsAffected() int64 {
	return delete.RowsAffected
}
func (delete *Delete) GetError() error {
	return delete.Error
}
func (delete *Delete) SetRawMap(rawJsonMap map[string]interface{}) {
	delete.RawMap = rawJsonMap
}
func (delete *Delete) SetCondition(condition map[string]interface{}) {
	delete.Condition = condition
}
func (delete *Delete) SetRowsAffected(rowsAffected int64) {
	delete.RowsAffected = rowsAffected
}
func (delete *Delete) SetRow(row interface{}) {
	delete.Row = row
}
func (delete *Delete) SetError(err error) {
	delete.Error = err
}

type TimeSummary struct {
	TodayAmount     uint64 `gorm:"column:today;type:int(11)"`
	YesterdayAmount uint64 `gorm:"column:yesterday;type:int(11)"`
	WeekAmount      uint64 `gorm:"column:week;type:int(11)"`
	MonthAmount     uint64 `gorm:"column:month;type:int(11)"`
}

type DBStructInterface interface {
	TableName() string
}

type SumData struct {
	SumFloat64 float64 `gorm:"column:sum_float64"`
}
