package rank_util_api

import (
	"app/services/request_out/download_center"
	"app/settings"
	"app/structs/check_include_util"
	"encoding/json"
	"github.com/panwenbin/ghttpclient"
	"strings"
)

type RankUtilApi struct {
	BaseUrl string
}

const (
	POST_REQUEST_BUILDER           = "/request-builder/:engine"
	POST_DOMAIN_INCLUDE_EXTRACTOR  = "/domain-include-extractor/:engine"
	POST_KEYWORD_INCLUDE_EXTRACTOR = "/keyword-include-extractor/:engine"
)

// NewRankUtilApi: 创建新的排名工具api
func NewRankUtilApi() *RankUtilApi {
	rankUtilApi := &RankUtilApi{}
	rankUtilApi.BaseUrl = settings.RankUtilApi

	return rankUtilApi
}

// apiUrl: 填充API参数并返回完整API地址
func (r *RankUtilApi) apiUrl(path string, params map[string]string) string {
	for key, value := range params {
		path = strings.Replace(path, key, value, 1)
	}

	return r.BaseUrl + path
}

// PostRequestBuilder: 请求构建DcRequest对象
func (r *RankUtilApi) PostRequestBuilder(engine, searchWord string, page int, capture bool, searchCycle int, priority string) (*download_center.DcRequest, error) {
	apiUrl := r.apiUrl(POST_REQUEST_BUILDER, map[string]string{
		":engine": engine,
	})

	requestBuilderRequest := check_include_util.RequestBuilderRequest{
		SearchWord:  searchWord,
		Page:        page,
		Capture:     capture,
		SearchCycle: searchCycle,
		Priority:    priority,
	}
	jsonBytes, err := json.Marshal(requestBuilderRequest)
	var dcRequest download_center.DcRequest
	err = ghttpclient.PostJson(apiUrl, jsonBytes, nil).ReadJsonClose(&dcRequest)
	if err != nil {
		return nil, err
	}

	return &dcRequest, nil
}

// PostIncludeExtractor: 收录解析
func (r *RankUtilApi) PostIncludeExtractor(html string, engine string) (*check_include_util.ParseIncludeResponse, error) {
	apiUrl := r.apiUrl(POST_DOMAIN_INCLUDE_EXTRACTOR, map[string]string{
		":engine": check_include_util.MapEngine[engine],
	})
	rankExtractor := check_include_util.ParseIncludeRequest{Body: html}
	jsonBytes, err := json.Marshal(rankExtractor)
	includeExtractorResponse := check_include_util.ParseIncludeResponse{}
	err = ghttpclient.PostJson(apiUrl, jsonBytes, nil).ReadJsonClose(&includeExtractorResponse)
	if err != nil {
		return nil, err
	}

	return &includeExtractorResponse, nil
}

func (r *RankUtilApi) PostKeywordIncludeExtractor(html string, engine string) (*check_include_util.KeywordParseIncludeResponse, error) {
	apiUrl := r.apiUrl(POST_KEYWORD_INCLUDE_EXTRACTOR, map[string]string{
		":engine": check_include_util.MapEngine[engine],
	})
	rankExtractor := check_include_util.ParseIncludeRequest{Body: html}
	jsonBytes, err := json.Marshal(rankExtractor)
	includeExtractorResponse := check_include_util.KeywordParseIncludeResponse{}
	err = ghttpclient.PostJson(apiUrl, jsonBytes, nil).ReadJsonClose(&includeExtractorResponse)
	if err != nil {
		return nil, err
	}

	return &includeExtractorResponse, nil
}
