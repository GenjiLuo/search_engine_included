package rank_task_service

import (
	"app/settings"
	"encoding/json"
	"errors"
	"github.com/panwenbin/ghttpclient"
	"net/http"
	"strings"
)

const PUT_分组发送不截图查排名关键词 = "/keywords/:check-match/:engine"

// apiUrl: 填充API参数并返回完整API地址
func apiUrl(path string, params map[string]string) string {
	baseUrl := settings.RankTaskApiBaseUrl

	for key, value := range params {
		path = strings.Replace(path, key, value, 1)
	}

	return baseUrl + path
}

func KeywordsPut(domain string, engine string, keywords []string) error {
	apiUrl := apiUrl(PUT_分组发送不截图查排名关键词, map[string]string{
		":check-match": domain,
		":engine":      engine,
	})

	putData, err := json.Marshal(keywords)
	if err != nil {
		return err
	}

	type keywordsPutResponse struct {
		Msg string `json:"msg"`
	}
	var result keywordsPutResponse
	client := ghttpclient.PutJson(apiUrl, putData, nil)

	err = client.ReadJsonClose(&result)
	if err != nil {
		return err
	}

	res, _ := client.Response()
	if res.StatusCode != http.StatusOK {
		return errors.New(result.Msg)
	}

	return nil
}
