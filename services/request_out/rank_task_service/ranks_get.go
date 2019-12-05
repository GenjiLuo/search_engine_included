package rank_task_service

import (
	"app/settings"
	"encoding/json"
	"github.com/panwenbin/ghttpclient"
	"github.com/panwenbin/ghttpclient/header"
	"strings"
)

const RANKS_GET = "/ranks/get/:check-match/:engine/:request-hash"

type RankResult struct {
	Word string `json:"word"`
	Rank int    `json:"rank"`
}

func getFullPath(api string) string {
	apiBaseUrl := settings.RankTaskApiBaseUrl
	return apiBaseUrl + api
}

func replace(oriString string, replaceItems map[string]string) string {
	result := oriString
	for oldString, newString := range replaceItems {
		result = strings.Replace(result, oldString, newString, 1)
	}
	return result
}

func RanksGet(domain string, engine string, requestHash string, keywords []string) ([]RankResult, error) {
	apiUrl := replace(getFullPath(RANKS_GET),
		map[string]string{
			":check-match":  domain,
			":engine":       engine,
			":request-hash": requestHash,
		})
	headers := header.GHttpHeader{}
	postData, err := json.Marshal(keywords)
	if err != nil {
		return nil, err
	}
	result := make([]RankResult, 0)
	err = ghttpclient.PostJson(apiUrl, postData, headers).ReadJsonClose(&result)
	return result, err
}
