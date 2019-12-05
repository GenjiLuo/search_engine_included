package actions

import (
	"app/services/keyword_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func KeywordsAddPut(c *gin.Context) {
	domain := c.Param("domain")
	keywords := make([]string, 0)
	err := c.BindJSON(&keywords)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求格式不正确",
		})
		return
	}

	err = keyword_service.KeywordsAdd(domain, keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "批量插入错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "数据保存成功",
	})
}
