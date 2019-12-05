package main

import (
	"app/databases"
	"app/databases/entities"
	"app/services"
	"app/structs/models/logics"
	"time"
)

func main() {
	finishedTask := &entities.KeywordIncludeTask{}
	databases.Db.Model(&entities.KeywordIncludeTask{}).Preload("KeywordInclude").Where("id=1").Find(&finishedTask)

	isIncluded := false
	keywordInclude := finishedTask.KeywordInclude
	if isIncluded == false {
		days := services.DaysApartToday(keywordInclude.LastIncludedAt)
		if days >= 1 {
			keywordInclude.NoIncludedDays = days
		}

		databases.Db.Model(&finishedTask).Update(entities.KeywordIncludeTask{Status: logics.TASK_STATUS_不被收录})
		if keywordInclude.NoIncludedDays >= logics.INCLUDE_调至不收录的连续无收录天数 {
			keywordInclude.IsIncluded = false
		}
	} else {
		keywordInclude.IsIncluded = true
		keywordInclude.NoIncludedDays = 0
		keywordInclude.LastIncludedAt = time.Now()
		databases.Db.Model(&finishedTask).Update(entities.KeywordIncludeTask{Status: logics.TASK_STATUS_被收录})
	}
	databases.Db.Save(&keywordInclude)
}
