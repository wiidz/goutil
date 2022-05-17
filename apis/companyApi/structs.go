package companyApi

type Status interface {
	bool | string
}

type RespInterface interface {
	// GetStatus() Status
	// GetStatus[T bool | string]() T
	GetStatus() interface{}
	GetErrorCode() int
	GetErrorMessage() string
	GetData() interface{}
	GetMessage() string
}

type CommonResp struct {
	//Status       T           `json:"status"` // 大部分都是true、false，但是模糊查询是200和201（没查到）,false
	Status       interface{} `json:"status"` // 大部分都是true、false，但是模糊查询是200和201（没查到）,false
	Message      string      `json:"message"`
	ErrorCode    int         `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

func (resp *CommonResp) GetStatus() interface{} {
	return resp.Status
}
func (resp *CommonResp) GetMessage() string {
	return resp.Message
}
func (resp *CommonResp) GetErrorCode() int {
	return resp.ErrorCode
}
func (resp *CommonResp) GetErrorMessage() string {
	return resp.ErrorMessage
}
func (resp *CommonResp) GetData() interface{} {
	return resp.Data
}

// FuzzyQueryParam 模糊查询企业
type FuzzyQueryParam struct {
	Name   string `json:"-"`
	PageNo int    `json:"pageNo"`
}
type FuzzyQueryResp struct {
	*CommonResp
	Data *FuzzyQueryData `json:"data"`
}
type FuzzyQueryData struct {
	Total int             `json:"total"`
	Num   int             `json:"num"`
	List  []*FuzzyCompany `json:"list"`
}
type FuzzyCompany struct {
	Name            string `json:"name"`
	LegalPersonName string `json:"legal_person_name"`
	RegCapital      string `json:"reg_capital"`
	RegDate         string `json:"reg_date"`
}

type QueryDetailParam struct {
	CompanyNameOrCreditNo string `json:"-"`
	IsRaiseErrorCode      int    `json:"isRaiseErrorCode,omitempty"` // 当请求传入不存在企业名称时是否抛出404错误。0为否，1为是，默认为否。可以避免传入不存在企业时扣减次数。
}

type QueryDetailResp struct {
	*CommonResp
	Data *QueryDetailData `json:"data"`
}

type QueryDetailData struct {
	StartDate       string `json:"startDate"`
	RegisterCapital string `json:"registerCapital"`
	Name            string `json:"name"`
	RegisterData    struct {
		Status        string `json:"status"`
		CreditNo      string `json:"creditNo"`
		OrgNo         string `json:"orgNo"`
		BusinessTerm  string `json:"businessTerm"`
		BelongOrg     string `json:"belongOrg"`
		RegType       string `json:"regType"`
		RegisterNo    string `json:"registerNo"`
		Address       string `json:"address"`
		BusinessScope string `json:"businessScope"`
	} `json:"registerData"`
	PartnerData struct {
		List []struct {
			TotalRealCapital   string `json:"totalRealCapital"`
			PartnerType        string `json:"partnerType"`
			TotalShouldCapital string `json:"totalShouldCapital"`
			PartnerName        string `json:"partnerName"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"partnerData"`
	ChangeRecordData struct {
		List []struct {
			Date   string `json:"date"`
			Item   string `json:"item"`
			After  string `json:"after"`
			Before string `json:"before"`
		} `json:"list"`
		HasMore bool `json:"hasMore"`
	} `json:"changeRecordData"`
	EmployeeData struct {
		List []struct {
			Name  string `json:"name"`
			Title string `json:"title"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"employeeData"`
	LegalPersonName string `json:"legalPersonName"`
}
