package aliCompanyApi

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const Domain = "https://api.81api.com/"

//const FuzzyQueryURL = "https://api.81api.com/fuzzyQueryCompanyInfo/" // "https://api.81api.com/fuzzyQueryCompanyInfo/[CompanyName]/"
//const QueryDetailURL = "https://api.81api.com/getCompanyBaseInfo/"   // http(s)://api.81api.com/getCompanyBaseInfo/[CompanyNameOrCreditNo]/

// CompanyApi 企业工商数据查询【天眼+启信】-查老板-查工商投融资-查专利商标-查工商年报-查工商风险失信-查法院公告-精准、模糊工商查询
// https://market.aliyun.com/products/56928005/cmapi029030.html?spm=5176.730005.result.50.a11f3cc6KxbxuD&innerSource=search_%E5%B7%A5%E5%95%86#sku=yuncode2303000004
// 15000次 = 1100元
type CompanyApi struct {
	Config *configStruct.AliApiConfig
}

// NewCompanyApi 企业工商数据查询接口
func NewCompanyApi(config *configStruct.AliApiConfig) *CompanyApi {
	return &CompanyApi{
		Config: config,
	}
}

// request 发送请求
// path 是不全的，要加上domain
func (api *CompanyApi) request(method networkStruct.Method, path string, params interface{}, iStruct RespInterface) (data interface{}, err error) {

	//【1】构建参数
	paramStr, _ := typeHelper.JsonEncode(params)
	paramMap := typeHelper.JsonDecodeMap(paramStr)

	//【3】发送请求
	res, err := networkHelper.MyRequest(&networkStruct.MyRequestParams{
		Method:      method,
		URL:         Domain + path,
		ContentType: networkStruct.BodyForm,
		Headers: map[string]string{
			"Authorization": "APPCODE " + api.Config.AppCode,
		},
		Params:    paramMap,
		ResStruct: iStruct,
	})

	//【4】判断结果
	if err != nil {
		return
	} else if res.StatusCode != 200 {
		err = errors.New("请求失败")
		return
	} else if !res.IsParsedSuccess {
		// 解析失败
		var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
		temp := NoDataResp{}
		err = json2.Unmarshal([]byte(res.ResStr), &temp)
		if err != nil {
			return
		}
		err = errors.New(temp.Data)
		return
	}

	if boolStatus, ok := iStruct.GetStatus().(bool); ok {
		if !boolStatus {
			err = errors.New(iStruct.GetErrorMessage())
		}
	} else if strStatus, ok := iStruct.GetStatus().(string); ok {
		if strStatus != "200" {
			err = errors.New(iStruct.GetErrorMessage())
		}
	}

	return
}

// request2 另一个发送请求（模糊查询专用）
// path 是不全的，要加上domain
func (api *CompanyApi) request2(method networkStruct.Method, path string, params interface{}, iStruct RespInterface) (data interface{}, err error) {

	//【1】构建参数
	paramStr, _ := typeHelper.JsonEncode(params)
	paramMap := typeHelper.JsonDecodeMap(paramStr)

	//【3】发送请求
	var statusCode int
	data, _, statusCode, err = networkHelper.RequestWithStructTest(method, networkStruct.BodyForm, Domain+path, paramMap, map[string]string{
		"Authorization": "APPCODE " + api.Config.AppCode,
	}, iStruct)

	//【4】判断结果
	if err != nil {
		return
	} else if statusCode != 200 {
		err = errors.New("请求失败")
	} else if iStruct.GetStatus() != true {
		err = errors.New(iStruct.GetErrorMessage())
	}

	return
}

// FuzzyQuery 企业工商数据模糊查询
func (api *CompanyApi) FuzzyQuery(params *FuzzyQueryParam) (*FuzzyQueryData, error) {
	res, err := api.request(networkStruct.Get, "fuzzyQueryCompanyInfo/"+params.Name+"/", params, &FuzzyQueryResp{})
	if err != nil {
		return nil, err
	}
	return res.(*FuzzyQueryResp).Data, err
}

// QueryDetail 企业工商数据精准查询
func (api *CompanyApi) QueryDetail(params *QueryDetailParam) (*QueryDetailData, error) {
	res, err := api.request(networkStruct.Get, "getCompanyBaseInfo/"+params.CompanyNameOrCreditNo+"/", params, &QueryDetailResp{})
	if err != nil {
		return nil, err
	}
	return res.(*QueryDetailResp).Data, err
}

// AbnormalInfo 企业经营异常信息
func (api *CompanyApi) AbnormalInfo(companyName string) (*AbnormalInfoData, error) {
	res, err := api.request(networkStruct.Get, "getCompanyAbnormalInfo/"+companyName+"/", nil, &AbnormalInfoResp{})
	if err != nil {
		return nil, err
	}
	return res.(*AbnormalInfoResp).Data, err
}

// LawsuitInfo 企业法律诉讼信息
func (api *CompanyApi) LawsuitInfo(companyName string) ([]*LawsuitInfoData, error) {
	res, err := api.request(networkStruct.Get, "getCompanyLawsuitInfo/"+companyName+"/", nil, &LawsuitInfoResp{})
	if err != nil {
		return nil, err
	}
	return res.(*LawsuitInfoResp).Data, err
}

// CourtInfo 企业法院公告信息
func (api *CompanyApi) CourtInfo(companyName string) ([]*CourtInfoData, error) {
	res, err := api.request(networkStruct.Get, "getCompanyCourtInfo/"+companyName+"/", nil, &CourtInfoResp{})
	if err != nil {
		return nil, err
	}
	return res.(*CourtInfoResp).Data, err
}
