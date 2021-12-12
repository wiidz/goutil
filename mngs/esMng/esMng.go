package esMng

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/wiidz/goutil/structs/configStruct"
	"log"
	"time"
)

var client *elastic.Client

type EsMng struct {
	client *elastic.Client
}

// Init 初始化
func Init(params *configStruct.EsConfig) (err error) {

	dsn := params.Host + ":" + params.Port
	log.Println("【es-dsn】", dsn)

	client, err = elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetSniff(false), elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetBasicAuth(params.Username, params.Password))

	return
}

// NewEsMng  获取es管理器
func NewEsMng() (es *EsMng) {
	return &EsMng{
		client: client,
	}
}

// Start 开启链接
func (es *EsMng) Start() {
	es.client.Start()
}

// Stop 关闭链接
func (es *EsMng) Stop() {
	es.client.Stop()
}

// Add 添加数据
func (es *EsMng) Add(index, id string, data interface{}) (res *elastic.IndexResponse,err error) {
	return es.client.Index().Index(index).Id(id).BodyJson(data).Do(context.Background())
}

// Update 修改
func (es *EsMng) Update(index, id string, data map[string]interface{}) ( res *elastic.UpdateResponse,err error) {
	return es.client.Update().Index(index).Id(id).Doc(data).Do(context.Background())
}

// DeleteByID 根据ID删除数据
func (es *EsMng) DeleteByID(index, id string) (res *elastic.DeleteResponse,err error) {
	return es.client.Delete().Index(index).Id(id).Do(context.Background())
}


// Truncate 清空一个index
func (es *EsMng) Truncate(index string) (res *elastic.IndicesDeleteResponse,err error) {
	return es.client.DeleteIndex(index).Do(context.Background())
}

// LikeQuery 多字段模糊查询
func (es *EsMng) LikeQuery(index string, page, pageSize int, searchStr string, searchFields ...string) (data []map[string]interface{}, err error) {

	//【1】拼接字符串
	matchPhraseQuery := elastic.NewMultiMatchQuery(searchStr, searchFields...)

	//【2】查询
	var res *elastic.SearchResult
	res, err = es.client.Search(index).Query(matchPhraseQuery).From(page * pageSize).Size(pageSize).Do(context.Background())
	if err != nil {
		return
	}
	if res.Hits.TotalHits.Value == 0 {
		return
	}

	//【3】构建数据
	for _, hit := range res.Hits.Hits {
		var tmp map[string]interface{}
		//err := json.Unmarshal(*hit.Source, &tmp)
		var bytes []byte
		bytes, err = hit.Source.MarshalJSON()
		err = json.Unmarshal(bytes, &tmp)
		if err != nil {
			return
		}
		data = append(data, tmp)
	}

	return
}

// Query 模糊查询数据
func (es *EsMng) Query(index, searchKey, searchStr string) (data []map[string]interface{}, err error) {

	//【1】拼接字符串
	matchPhraseQuery1 := elastic.NewMatchPhraseQuery(searchKey, searchStr)
	fmt.Println("es search:", searchKey, searchStr)
	fmt.Println(matchPhraseQuery1)

	//【2】查询
	var res *elastic.SearchResult
	res, err = es.client.Search(index).Query(matchPhraseQuery1).From(0).Size(300).Do(context.Background())
	if err != nil {
		return
	}
	if res.Hits.TotalHits.Value == 0 {
		return
	}

	//【3】处理数据
	for _, hit := range res.Hits.Hits {
		var tmp map[string]interface{}
		//err := json.Unmarshal(*hit.Source, &tmp)
		var bytes []byte
		bytes, err = hit.Source.MarshalJSON()
		err = json.Unmarshal(bytes, &tmp)
		if err != nil {
			return
		}
		data = append(data, tmp)
	}

	return
}

// QueryByField 指定字段模糊查询数据
func (es *EsMng) QueryByField(index, searchKey, searchStr, field string) (data []interface{}, err error) {

	//【1】查询
	var res *elastic.SearchResult
	matchPhraseQuery1 := elastic.NewMatchPhraseQuery(searchKey, searchStr)
	res, err = es.client.Search(index).Query(matchPhraseQuery1).From(0).Size(300).Do(context.Background())
	if err != nil {
		return
	}
	if res.Hits.TotalHits.Value == 0 {
		return
	}

	//【2】处理数据
	for _, hit := range res.Hits.Hits {
		var tmp map[string]interface{}
		//err := json.Unmarshal(*hit.Source, &tmp)
		var bytes []byte
		bytes, err = hit.Source.MarshalJSON()
		err = json.Unmarshal(bytes, &tmp)
		if err != nil {
			return
		}
		data = append(data, tmp[field])
	}
	return
}
