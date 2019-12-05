package actions

import (
	"app/services/keyword_service"
	"app/services/request_out/rank_task_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RanksGet(c *gin.Context) {
	checkMatch := c.Param("check-match")
	engine := c.Param("engine")
	requestHash := c.Param("request-hash")

	filteredKeywords, err := keyword_service.IncludedKeywords(checkMatch, engine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "过滤关键词失败",
		})
		return
	}
	result, err := rank_task_service.RanksGet(checkMatch, engine, requestHash, filteredKeywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "查排名接口调用失败",
		})
		return
	}
	c.JSON(http.StatusOK, result)
}
