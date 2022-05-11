package imMsg

import "github.com/wiidz/goutil/mngs/yunxinMng/imClient"

type DeleteQuestParam struct {
	StartTime   int64  `json:"startTime" validate:"required"` // 被清理文件的开始时间(毫秒级)
	EndTime     int64  `json:"endTime" validate:"required"`   // 被清理文件的结束时间(毫秒级)
	ContentType string `json:"contentType"`                   // 被清理的文件类型，文件类型包含contentType则被清理 如原始文件类型为"image/png"，contentType参数为"image",则满足被清理条件
	Tag         string `json:"tag"`                           // 被清理文件的应用场景，完全相同才被清理 如上传文件时知道场景为"usericon",tag参数为"usericon"，则满足被清理条件
}

type DeleteQuestResp struct {
	*imClient.CommonResp
	Data struct {
		TaskID string `json:"taskid"`
	} `json:"data"`
}

// DeleteQuest 上传NOS文件清理任务
// https://doc.yunxin.163.com/docs/TM5MzM5Njk/DEwMTE3NzQ?platformId=60353#上传NOS文件清理任务
// 上传NOS文件清理任务，按时间范围和文件类下、场景清理符合条件的文件
// 每天提交的任务数量有限制，请合理规划
// 关于startTime与endTime请注意：
// startTime必须小于endTime且大于0，endTime和startTime差值在1天以上，7天以内。
// endTime必须早于今天（即只可以清理今天以前的文件）。
func (api *Api) DeleteQuest(param *DeleteQuestParam) (*DeleteQuestResp, error) {
	res, err := api.Client.Post(SubDomain+"del.action", param, &DeleteQuestResp{})
	return res.(*DeleteQuestResp), err
}
