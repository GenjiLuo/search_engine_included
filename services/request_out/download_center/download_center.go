package download_center

import (
	"app/settings"
	"encoding/json"
	"errors"
	"github.com/panwenbin/ghttpclient"
	"github.com/panwenbin/ghttpclient/header"
	"strings"
)

type DownloadCenter struct {
	BaseUrl string
}

const PUT_REQUEST = "/requests/:unique-key"
const GET_REQUESTS = "/requests"
const POST_RESPONSES_CHECK = "/responses-check"
const GET_RESPONSE = "/responses/:unique-key"
const PUT_RESET_REQUEST = "/request-reset/:unique-key"

// NewDownloadCenter: 创建一个DownloadCenter对象
func NewDownloadCenter() *DownloadCenter {
	downloadCenter := &DownloadCenter{}
	downloadCenter.BaseUrl = settings.DcWrapperApi

	return downloadCenter
}

// apiUrl: 填充API参数并返回完整API地址
func (dc *DownloadCenter) apiUrl(path string, params map[string]string) string {
	for key, value := range params {
		path = strings.Replace(path, key, value, 1)
	}

	return dc.BaseUrl + path
}

// PutRequest: 发送单个任务到下载中心
func (dc *DownloadCenter) PutRequest(dcRequest *DcRequest) error {
	apiUrl := dc.apiUrl(PUT_REQUEST, map[string]string{
		":unique-key": dcRequest.UniqueKey,
	})

	jsonBytes, err := json.Marshal(dcRequest)
	client := ghttpclient.PutJson(apiUrl, jsonBytes, nil)

	var dcRequestRepeat DcRequest
	err = client.ReadJsonClose(&dcRequestRepeat)
	if err != nil {
		return err
	}

	res, _ := client.Response()
	if res.StatusCode != 200 {
		return errors.New("status code != 200")
	}
	if dcRequestRepeat.UniqueKey != dcRequest.UniqueKey {
		return errors.New("unique key not match")
	}

	return nil
}

// PostResponsesCheck: 批量检查下载中心任务完成情况
func (dc *DownloadCenter) PostResponsesCheck(uniqueKeys []string) ([]string, error) {
	apiUrl := dc.apiUrl(POST_RESPONSES_CHECK, nil)
	jsonBytes, _ := json.Marshal(uniqueKeys)
	headers := make(header.GHttpHeader)
	headers.Set("Accept-Encoding", "gzip")

	var finishedUniqueKeys []string
	err := ghttpclient.PostJson(apiUrl, jsonBytes, headers).ReadJsonClose(&finishedUniqueKeys)
	if err != nil {
		return nil, err
	}

	return finishedUniqueKeys, nil
}

// GetResponse: 获取某个任务的结果
func (dc *DownloadCenter) GetResponse(uniqueKey string) (*DcResponse, error) {
	apiUrl := dc.apiUrl(GET_RESPONSE, map[string]string{
		":unique-key": uniqueKey,
	})
	headers := make(header.GHttpHeader)
	headers.AcceptEncodingGzip()

	var dcResponse DcResponse
	err := ghttpclient.Get(apiUrl, headers).ReadJsonClose(&dcResponse)
	if err != nil {
		return nil, err
	}

	return &dcResponse, nil
}

// ResetRequest: 重置某个任务
func (dc *DownloadCenter) ResetRequest(uniqueKey string) error {
	apiUrl := dc.apiUrl(PUT_RESET_REQUEST, map[string]string{
		":unique-key": uniqueKey,
	})
	headers := make(header.GHttpHeader)
	headers.AcceptEncodingGzip()
	type response struct {
		Msg string `json:"msg"`
	}
	res := response{}

	err := ghttpclient.Get(apiUrl, headers).ReadJsonClose(&res)
	if err != nil {
		return err
	}

	return nil
}
