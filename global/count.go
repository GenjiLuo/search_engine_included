package global

import (
	"app/databases"
	"app/databases/entities"
	"app/databases/scopes/task_scope"
)

var BeforeQueriedDomainCount int64
var BeforeQueriedKeywordCount int64

func init() {
	ReadBeforeQueriedCount()
}

func ReadBeforeQueriedCount() {
	databases.Db.Model(&entities.DomainIncludeTask{}).Scopes(task_scope.BeforeQueried).Count(&BeforeQueriedDomainCount)
	databases.Db.Model(&entities.KeywordIncludeTask{}).Scopes(task_scope.BeforeQueried).Count(&BeforeQueriedKeywordCount)
}
