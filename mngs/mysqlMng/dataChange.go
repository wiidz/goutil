package mysqlMng

import (
	"gorm.io/gorm"
)

type SteppingWay int8 // 步进方式

const Increase SteppingWay = 1 // 增
const Decrease SteppingWay = 2 // 减

type DataSteppingParam struct {
	ID        uint64
	Way       SteppingWay
	Model     interface{}
	FieldName string
	Amount    int64
}

func (p *DataSteppingParam) GetOperator() (operator string) {

	if p.Way == Increase {
		operator = " - "
	} else {
		operator = " + "
	}
	return
}

// DataStepping 针对一个数字字段递增或递减
func DataStepping(conn *gorm.DB, params *DataSteppingParam) (err error) {
	return conn.Model(&params.Model).Where("id = ?", &params.ID).Update(params.FieldName, gorm.Expr(params.FieldName+params.GetOperator()+" ?", params.Amount)).Error
}
