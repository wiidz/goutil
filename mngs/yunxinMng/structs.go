package yunxinMng

type MsgType int

const Text MsgType = 0  // 文本消息
const Img MsgType = 1   // 图片消息
const Audio MsgType = 2 // 语音消息
const Video MsgType = 3 // 视频消息
const Geo MsgType = 4   // 地理位置消息

const TeamNotice MsgType = 5     // 群通知
const ChatRoomNotice MsgType = 5 // 聊天室通知（和上面的type都是5）

const File MsgType = 6       // 文件消息
const Tips MsgType = 10      // 提示消息
const Customer MsgType = 100 // 自定义消息

type MsgInterface interface {
	Hello()
}

// TextMsg 文本消息
type TextMsg struct {
	Msg string `json:"msg"` // 消息内容
}

// ImgMsg 图片消息
type ImgMsg struct {
	Name string `json:"name"` // 图片name
	Md5  string `json:"md5"`  // 图片文件md5
	Url  string `json:"url"`  // 生成的url
	Ext  string `json:"ext"`  // 图片后缀
	W    int    `json:"w"`    // 宽
	H    int    `json:"h"`    // 高
	Size int    `json:"size"` // 图片文件大小
}

func (msg ImgMsg) Hello() {}

// AudioMsg 语音消息
type AudioMsg struct {
	Dur  int    `json:"dur"`  // 语音持续时长ms
	Md5  string `json:"md5"`  // 语音文件的md5值
	Url  string `json:"url"`  // 生成的url
	Ext  string `json:"ext"`  // 语音消息格式，只能是aac格式
	Size int    `json:"size"` // 语音文件大小
}

func (msg *AudioMsg) Hello() {}

// VideoMsg 视频消息
type VideoMsg struct {
	Dur  int    `json:"dur"`  // 视频持续时长ms
	Md5  string `json:"md5"`  // 视频文件的md5值
	Url  string `json:"url"`  // 生成的url
	W    int    `json:"w"`    // 宽
	H    int    `json:"h"`    // 高
	Ext  string `json:"ext"`  // 视频格式
	Size int    `json:"size"` // 视频文件大小
}

func (msg *VideoMsg) Hello() {}

// GeoMsg 位置消息
type GeoMsg struct {
	Title string  `json:"title"` // 地理位置title
	Lng   float64 `json:"lng"`   // 经度
	Lat   float64 `json:"lat"`   // 纬度
}

func (msg *GeoMsg) Hello() {}

// FileMsg 文件消息
type FileMsg struct {
	Name string `json:"name"` // 文件名
	Md5  string `json:"md5"`  // 文件MD5
	Url  string `json:"url"`  // 生成的url
	Ext  string `json:"ext"`  // 文件后缀类型
	Size int    `json:"size"` // 大小
}

func (msg *FileMsg) Hello() {}

// TipsMsg 提示消息
type TipsMsg struct {
	Msg string `json:"msg"` // 消息内容
}

func (msg *TipsMsg) Hello() {}

// CustomerMsg 第三方定义的body体，json格式
type CustomerMsg struct {
}

func (msg *CustomerMsg) Hello() {}

// MsgNoticeMsg 群通知
type MsgNoticeMsg struct {
	TeamID       int      `json:"tid"`          //群id
	TeamName     string   `json:"tname"`        //群名称 （某些操作会有）
	Ope          int      `json:"ope"`          // notify通知类型 （0:群拉人，1:群踢人，2:退出群，3:群信息更新，4:群解散，5:申请加入群成功，6:退出并移交群主，7:增加管理员，8:删除管理员，9:接受邀请进群）
	Accids       []string `json:"accids"`       // 被操作的对象 （群成员操作时才有）
	Intro        string   `json:"intro"`        //（ope=3时群信息修改项）
	Announcement string   `json:"announcement"` //（ope=3时群信息修改项）
	JoinMode     int      `json:"joinmode"`     //加入群的模式0不需要认证，1需要认证 ,（ope=3时群信息修改项）
	Config       string   `json:"config"`       //（ope=3时群信息修改项）
	UpdateTime   int64    `json:"updatetime"`   //通知后台更新时间 （群操作时有）
}

func (msg *MsgNoticeMsg) Hello() {}

// ChatRoomNoticeMsg 聊天室通知
type ChatRoomNoticeMsg struct {
	ID   int `json:"id"` // //301: 成员进入聊天室,302: 成员离开聊天室,303: 成员被加黑,304: 成员被取消黑名单,305: 成员被设置禁言,306: 成员被取消禁言,307: 设置为管理员,308: 取消管理员,309: 成员设定为固定成员,310: 成员取消固定成员,312: 聊天室信息更新,313: 成员被踢,314: 新增临时禁言,315: 主动解除临时禁言,316: 成员更新聊天室内的角色信息,317: 麦序队列中有变更,318: 聊天室禁言,319: 聊天室解除禁言状态,320: 麦序队列中有批量变更
	Data struct {
		OpeNick  string   `json:"opeNick"`
		Operator string   `json:"operator"`
		Target   []string `json:"target"`
		TarNick  []string `json:"tarNick"`
		Ext      string   `json:"ext"`
	} `json:"data"`
}

func (msg *ChatRoomNoticeMsg) Hello() {}
