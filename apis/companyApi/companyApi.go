package companyApi

import (
	"errors"
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
	var statusCode int
	data, _, statusCode, err = networkHelper.RequestWithStructTest(method, networkStruct.BodyForm, Domain+path, paramMap, map[string]string{
		"Authorization": "APPCODE " + api.Config.AppCode,
	}, iStruct)

	//【4】判断结果
	if err != nil {
		return
	} else if statusCode != 200 {
		err = errors.New("请求失败")
		return
	}

	if boolStatus, ok := iStruct.GetStatus().(bool); ok {
		if !boolStatus {
			err = errors.New(iStruct.GetMessage())
		}
	} else if strStatus, ok := iStruct.GetStatus().(string); ok {
		if strStatus != "200" {
			err = errors.New(iStruct.GetMessage())
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
		err = errors.New(iStruct.GetMessage())
	}

	return
}

// FuzzyQuery 企业工商数据模糊查询
func (api *CompanyApi) FuzzyQuery(params *FuzzyQueryParam) (*FuzzyQueryResp, error) {
	res, err := api.request(networkStruct.Get, "fuzzyQueryCompanyInfo/"+params.Name+"/", params, &FuzzyQueryResp{})
	return res.(*FuzzyQueryResp), err
}

// QueryDetail 企业工商数据精准查询
func (api *CompanyApi) QueryDetail(params *QueryDetailParam) (resp *QueryDetailResp, err error) {
	res, err := api.request(networkStruct.Get, "getCompanyBaseInfo/"+params.CompanyNameOrCreditNo+"/", params, &QueryDetailResp{})
	return res.(*QueryDetailResp), err
}
