package services

import (
	"app/databases"
	"app/databases/entities"
	"app/services/request_out/download_center"
	"app/services/request_out/rank_util_api"
	"app/structs/check_include_util"
	"errors"
	"fmt"
)

type ParserService struct {
}

// BuildDomainDcRequest: 构建单个下载中心任务结构
func (parserService *ParserService) BuildDomainDcRequest(task *entities.DomainIncludeTask) (*download_center.DcRequest, error) {
	var engined_domain entities.EnginedDomain
	databases.Db.Model(task).Association("engined_domain").Find(&engined_domain)

	engine, ok := check_include_util.MapEngine[engined_domain.Engine]
	if !(ok) {
		return &download_center.DcRequest{}, errors.New("BuildDomainDcRequest get engine error")
	}

	searchWord := "site:" + engined_domain.Domain
	rankUrlApi := rank_util_api.NewRankUtilApi()
	return rankUrlApi.PostRequestBuilder(engine, searchWord, 1, false, 0, "normal")
}

// BuildKeywordDcRequest: 构建单个下载中心任务结构
func (parserService *ParserService) BuildKeywordDcRequest(task *entities.KeywordIncludeTask) (*download_center.DcRequest, error) {
	var keyword_include entities.KeywordInclude
	var engined_domain entities.EnginedDomain
	databases.Db.Model(task).Association("KeywordInclude").Find(&keyword_include)
	databases.Db.Model(keyword_include).Association("EnginedDomain").Find(&engined_domain)

	engine, ok := check_include_util.MapEngine[engined_domain.Engine]
	if !(ok) {
		return &download_center.DcRequest{}, errors.New("BuildDomainDcRequest get engine error")
	}

	searchWord := fmt.Sprintf("site:%s %s", engined_domain.Domain, keyword_include.Keyword)
	rankUrlApi := rank_util_api.NewRankUtilApi()
	return rankUrlApi.PostRequestBuilder(engine, searchWord, 1, false, 0, "normal")
}

// ParseInclude: 解析某平台收录
func (parserService *ParserService) ParseInclude(html string, engine string) (int, error) {
	rankUrlApi := rank_util_api.NewRankUtilApi()
	res, err := rankUrlApi.PostIncludeExtractor(html, engine)
	if err != nil {
		return 0, err
	}

	return res.IncludeNum, nil
}

// KeywordParseInclude: 解析某关键词收录
func (parserService *ParserService) KeywordParseInclude(html string, engine string) (bool, error) {
	rankUrlApi := rank_util_api.NewRankUtilApi()
	res, err := rankUrlApi.PostKeywordIncludeExtractor(html, engine)
	if err != nil {
		return false, err
	}

	return res.IsIncluded, nil
}
