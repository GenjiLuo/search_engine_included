package routers

import (
	"app/actions"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Load() *gin.Engine {
	r.PUT("/keywords-add/:domain", actions.KeywordsAddPut)
	r.PUT("/keywords-remove/:domain", actions.KeywordsRemove)
	r.GET("/ranks/:check-match/:engine/:request-hash", actions.RanksGet)
	r.PUT("/rank-task/:domain", actions.RankTaskPut)

	return r
}

func init() {
	r = gin.Default()
}
