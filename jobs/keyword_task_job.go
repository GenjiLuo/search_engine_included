package jobs

import (
	"app/channels"
	"app/databases"
	"app/databases/db_keyword_include_service"
	"app/databases/entities"
	"app/databases/scopes/task_scope"
	"app/global"
	"app/services"
	"app/services/keyword_service"
	"app/services/request_out/download_center"
	"app/services/request_out/rank_task_service"
	"app/services/task_service"
	"app/structs/models/logics"
	"sync"
	"sync/atomic"
	"time"
)

func AddKeywordsToTask() {
	if global.BeforeQueriedKeywordCount < logics.TASK_查询完成前任务数量限制 {
		keywordsCount := services.AddSomeUnCheckKeywordToTasks()
		if keywordsCount > 0 {
			atomic.AddInt64(&global.BeforeQueriedKeywordCount, int64(keywordsCount))
			return
		}
	}
	time.Sleep(time.Second * 3)
}

// SendKeywordTasksToChan: 从Chan发送任务
func SendKeywordTasksToChan() {
	var keywordTasks []*entities.KeywordIncludeTask
	limit := cap(channels.KeywordTaskSendingChan) - len(channels.KeywordTaskSendingChan)
	if limit == 0 {
		time.Sleep(time.Second)
		return
	}
	databases.Db.Model(&entities.KeywordIncludeTask{}).Scopes(task_scope.UnQueried).Limit(limit).Scan(&keywordTasks)

	if len(keywordTasks) != 0 {
		task_service.SendKeywordTasksToChan(keywordTasks, channels.KeywordTaskSendingChan, func(keywordTask *entities.KeywordIncludeTask) {
			keywordTask.Status = logics.TASK_STATUS_查询中
			databases.Db.Save(&keywordTask)
		})
	} else {
		time.Sleep(time.Second * 5)
	}
}

// SendKeywordDcRequestsFromChan: 发送下载中心
func SendKeywordDcRequestsFromChan() {
	task := <-channels.KeywordTaskSendingChan
	parserService := &services.ParserService{}
	dcRequest, err := parserService.BuildKeywordDcRequest(task)
	if err != nil {
		time.Sleep(time.Second * 5)
		db_keyword_include_service.NeedCheck([]int{task.KeywordIncludeId})
		databases.Db.Model(&task).Update(entities.KeywordIncludeTask{Status: logics.TASK_STATUS_未查询})
		return
	}

	dc := download_center.NewDownloadCenter()
	err = dc.PutRequest(dcRequest)
	if err != nil {
		time.Sleep(time.Second * 5)
		db_keyword_include_service.NeedCheck([]int{task.KeywordIncludeId})
		databases.Db.Model(&task).Update(entities.KeywordIncludeTask{Status: logics.TASK_STATUS_未查询})
		return
	}

	task.UniqueKey = dcRequest.UniqueKey
	databases.Db.Save(&task)
}

// FetchKeywordQueryParsedResult: 获取查询和解析结果
func FetchKeywordQueryParsedResult() {
	var beforeQueriedKeywordTasks []*entities.KeywordIncludeTask
	databases.Db.Preload("KeywordInclude").Scopes(task_scope.Querying).Find(&beforeQueriedKeywordTasks)
	if len(beforeQueriedKeywordTasks) == 0 {
		time.Sleep(time.Second * 2)
		return
	}

	fetchKeywordTasksResults(beforeQueriedKeywordTasks)
	time.Sleep(time.Second)
}

type KeywordUniqueKeyTaskGroup struct {
	UniqueKey string
	Tasks     []*entities.KeywordIncludeTask
}

func fetchKeywordTasksResults(tasks []*entities.KeywordIncludeTask) {
	var uniqueKeys []string
	uniqueKeyIndexedTasks := make(map[string]*entities.KeywordIncludeTask)
	for _, task := range tasks {
		uniqueKeys = append(uniqueKeys, task.UniqueKey)
		uniqueKeyIndexedTasks[task.UniqueKey] = task
	}

	dc := download_center.NewDownloadCenter()
	finishedUniqueKeys, err := dc.PostResponsesCheck(uniqueKeys)
	if err != nil {
		time.Sleep(time.Second * 5)
		return
	}
	if len(finishedUniqueKeys) == 0 {
		time.Sleep(time.Second * 30)
		return
	}

	var finishedTasks []*entities.KeywordIncludeTask
	databases.Db.
		Preload("KeywordInclude").
		Preload("KeywordInclude.EnginedDomain").
		Scopes(task_scope.UniqueKeysIn(finishedUniqueKeys)).
		Find(&finishedTasks)

	uniqueKeyTasksMap := task_service.KeywordUniqueKeyMappedTasks(finishedTasks)

	uniqueKeyTaskGroupChan := make(chan KeywordUniqueKeyTaskGroup)
	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			for {
				uniqueKeyTaskGroup, ok := <-uniqueKeyTaskGroupChan
				if !ok {
					break
				}
				tryFinishKeywordTask(uniqueKeyTaskGroup)
			}
			wg.Done()
		}()
	}
	for _, uniqueKey := range finishedUniqueKeys {
		uniqueKeyTaskGroup := KeywordUniqueKeyTaskGroup{
			UniqueKey: uniqueKey,
			Tasks:     uniqueKeyTasksMap[uniqueKey],
		}
		uniqueKeyTaskGroupChan <- uniqueKeyTaskGroup
		time.Sleep(time.Microsecond * 10)
	}
	close(uniqueKeyTaskGroupChan)

	wg.Wait()
}

// tryFinishKeywordTask: 试图结束任务
func tryFinishKeywordTask(group KeywordUniqueKeyTaskGroup) {
	dc := download_center.NewDownloadCenter()
	dcResponse, err := dc.GetResponse(group.UniqueKey)
	if err != nil {
		taskCount := int64(len(group.Tasks))
		atomic.AddInt64(&global.BeforeQueriedKeywordCount, -taskCount)

		databases.Db.Model(&entities.KeywordIncludeTask{}).
			Where(entities.KeywordIncludeTask{UniqueKey: group.UniqueKey}).
			Updates(entities.KeywordIncludeTask{Status: logics.TASK_STATUS_查询失败})

		keywordIncludeId := make([]int, 0)
		databases.Db.Model(&entities.KeywordIncludeTask{}).
			Where(entities.KeywordIncludeTask{UniqueKey: group.UniqueKey}).
			Pluck("keyword_include_id", &keywordIncludeId)
		db_keyword_include_service.NeedCheck(keywordIncludeId)

		return
	}
	if dcResponse.Body == "" {
		_ = download_center.NewDownloadCenter().ResetRequest(group.UniqueKey)
		atomic.AddInt64(&global.BeforeQueriedKeywordCount, 1)
		return
	}

	for _, finishedTask := range group.Tasks {
		atomic.AddInt64(&global.BeforeQueriedKeywordCount, -1)
		parserService := services.ParserService{}
		isIncluded, err := parserService.KeywordParseInclude(dcResponse.Body, finishedTask.KeywordInclude.EnginedDomain.Engine)
		if err != nil {
			db_keyword_include_service.NeedCheck([]int{finishedTask.KeywordIncludeId})

			databases.Db.Model(&finishedTask).Updates(entities.KeywordIncludeTask{Status: logics.TASK_STATUS_查询失败})
			continue
		}

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
}

// SyncRankTasks 同步查排名任务
func SyncRankTasks() {
	enginedDomains := make([]entities.EnginedDomain, 0)
	databases.Db.Model(entities.EnginedDomain{}).
		Preload("KeywordIncludes", "is_included = ?", keyword_service.IS_INCLUDE_YES).
		Where("include_num > ?", 0).
		Find(&enginedDomains)

	for i := range enginedDomains {
		keywords := make([]string, 0)
		for j := range enginedDomains[i].KeywordIncludes {
			keywords = append(keywords, enginedDomains[i].KeywordIncludes[j].Keyword)
		}
		if len(keywords) == 0 {
			continue
		}
		_ = rank_task_service.KeywordsPut(enginedDomains[i].Domain, enginedDomains[i].Engine, keywords)
	}

	time.Sleep(time.Hour * 24)
}
