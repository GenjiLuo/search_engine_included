package actions

import (
	"app/databases"
	"app/databases/entities"
	"app/services/keyword_service"
	"app/services/request_out/rank_task_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RankTaskPut(c *gin.Context) {
	domain := c.Param("domain")

	enginedDomains := make([]entities.EnginedDomain, 0)
	databases.Db.Model(entities.EnginedDomain{}).
		Preload("KeywordIncludes", "is_included = ?", keyword_service.IS_INCLUDE_YES).
		Where("domain = ?", domain).
		Find(&enginedDomains)

	if len(enginedDomains) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "未找到该域名信息",
		})
	}

	for i := range enginedDomains {
		keywords := make([]string, 0)
		for j := range enginedDomains[i].KeywordIncludes {
			keywords = append(keywords, enginedDomains[i].KeywordIncludes[j].Keyword)
		}
		if len(keywords) == 0 {
			continue
		}
		err := rank_task_service.KeywordsPut(enginedDomains[i].Domain, enginedDomains[i].Engine, keywords)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
