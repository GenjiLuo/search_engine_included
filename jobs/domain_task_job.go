package jobs

import (
	"app/channels"
	"app/common/debug_log"
	"app/databases"
	"app/databases/entities"
	"app/databases/scopes/task_scope"
	"app/global"
	"app/services"
	"app/services/request_out/download_center"
	"app/services/task_service"
	"app/structs/models/logics"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func AddDomainsToTask() {
	debug_log.Info(fmt.Sprintf("BeforeQueriedTasksCount: %d", global.BeforeQueriedDomainCount), "COUNT")
	if global.BeforeQueriedDomainCount < logics.TASK_查询完成前任务数量限制 {
		keywordsCount := services.AddSomeUnCheckDomainToTasks()
		if keywordsCount > 0 {
			atomic.AddInt64(&global.BeforeQueriedDomainCount, int64(keywordsCount))
			return
		}
	}
	time.Sleep(time.Second * 3)
}

// SendDomainTasksToChan: 从Chan发送任务
func SendDomainTasksToChan() {
	var domainTasks []*entities.DomainIncludeTask
	limit := cap(channels.DomainTaskSendingChan) - len(channels.DomainTaskSendingChan)
	if limit == 0 {
		time.Sleep(time.Second)
		return
	}
	databases.Db.Model(&entities.DomainIncludeTask{}).Scopes(task_scope.UnQueried).Limit(limit).Scan(&domainTasks)

	if len(domainTasks) != 0 {
		task_service.SendDomainTasksToChan(domainTasks, channels.DomainTaskSendingChan, func(domainTask *entities.DomainIncludeTask) {
			domainTask.Status = logics.TASK_STATUS_查询中
			databases.Db.Save(&domainTask)
		})
	} else {
		time.Sleep(time.Second * 5)
	}
}

// SendDomainDcRequestsFromChan: 发送下载中心
func SendDomainDcRequestsFromChan() {
	task := <-channels.DomainTaskSendingChan
	parserService := &services.ParserService{}
	dcRequest, err := parserService.BuildDomainDcRequest(task)
	if err != nil {
		time.Sleep(time.Second * 5)
		databases.Db.Model(&task).Update(entities.DomainIncludeTask{Status: logics.TASK_STATUS_未查询})
		return
	}

	dc := download_center.NewDownloadCenter()
	err = dc.PutRequest(dcRequest)
	if err != nil {
		time.Sleep(time.Second * 5)
		databases.Db.Model(&task).Update(entities.DomainIncludeTask{Status: logics.TASK_STATUS_未查询})
		return
	}

	task.UniqueKey = dcRequest.UniqueKey
	databases.Db.Save(&task)
}

type UniqueKeyTaskGroup struct {
	UniqueKey string
	Tasks     []*entities.DomainIncludeTask
}

// FetchDomainQueryParsedResult: 获取查询和解析结果
func FetchDomainQueryParsedResult() {
	var beforeQueriedDomainTasks []*entities.DomainIncludeTask
	databases.Db.Preload("EnginedDomain").Scopes(task_scope.Querying).Find(&beforeQueriedDomainTasks)
	if len(beforeQueriedDomainTasks) == 0 {
		time.Sleep(time.Second * 2)
		return
	}

	fetchDomainTasksResults(beforeQueriedDomainTasks)
	time.Sleep(time.Second)
}

// fetchDomainTasksResults: 获取查询结果
func fetchDomainTasksResults(tasks []*entities.DomainIncludeTask) {
	var uniqueKeys []string
	uniqueKeyIndexedTasks := make(map[string]*entities.DomainIncludeTask)
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

	var finishedTasks []*entities.DomainIncludeTask
	databases.Db.Preload("EnginedDomain").Scopes(task_scope.UniqueKeysIn(finishedUniqueKeys)).Find(&finishedTasks)
	uniqueKeyTasksMap := task_service.UniqueKeyMappedTasks(finishedTasks)
	uniqueKeyTaskGroupChan := make(chan UniqueKeyTaskGroup)
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			for {
				uniqueKeyTaskGroup, ok := <-uniqueKeyTaskGroupChan
				if !ok {
					break
				}
				tryFinishDomainTask(uniqueKeyTaskGroup)
			}
			wg.Done()
		}()
	}

	for _, uniqueKey := range finishedUniqueKeys {
		uniqueKeyTaskGroup := UniqueKeyTaskGroup{
			UniqueKey: uniqueKey,
			Tasks:     uniqueKeyTasksMap[uniqueKey],
		}
		uniqueKeyTaskGroupChan <- uniqueKeyTaskGroup
		time.Sleep(time.Microsecond * 10)
	}
	close(uniqueKeyTaskGroupChan)

	wg.Wait()
}

// tryFinishDomainTask: 试图结束任务
func tryFinishDomainTask(group UniqueKeyTaskGroup) {
	dc := download_center.NewDownloadCenter()
	dcResponse, err := dc.GetResponse(group.UniqueKey)
	if err != nil {
		taskCount := int64(len(group.Tasks))
		atomic.AddInt64(&global.BeforeQueriedDomainCount, -taskCount)
		databases.Db.Model(&entities.DomainIncludeTask{}).Where(entities.DomainIncludeTask{UniqueKey: group.UniqueKey}).Updates(entities.DomainIncludeTask{Status: logics.TASK_STATUS_查询失败})
		return
	}
	if dcResponse.Body == "" {
		_ = download_center.NewDownloadCenter().ResetRequest(group.UniqueKey)
		atomic.AddInt64(&global.BeforeQueriedDomainCount, 1)
		return
	}

	for _, finishedTask := range group.Tasks {
		atomic.AddInt64(&global.BeforeQueriedDomainCount, -1)
		databases.Db.Model(&finishedTask.EnginedDomain).Update(&entities.EnginedDomain{CheckedAt: time.Now()})

		parserService := services.ParserService{}
		includeNum, err := parserService.ParseInclude(dcResponse.Body, finishedTask.EnginedDomain.Engine)
		if err != nil {
			databases.Db.Model(&finishedTask).Updates(entities.DomainIncludeTask{Status: logics.TASK_STATUS_查询失败})
			continue
		}

		oldIncludeNum := finishedTask.EnginedDomain.IncludeNum
		databases.Db.Model(&finishedTask.EnginedDomain).Update(map[string]interface{}{"include_num": includeNum})

		if includeNum == logics.INCLUDE_无收录 {
			databases.Db.Model(&finishedTask).Update(entities.DomainIncludeTask{Status: logics.TASK_STATUS_不被收录})
			databases.Db.Model(&entities.KeywordInclude{}).Where("engined_domain_id = ?", finishedTask.EnginedDomain.ID).Update(map[string]interface{}{
				"is_included": false,
			})
			continue
		}

		databases.Db.Model(&finishedTask).Update(entities.DomainIncludeTask{Status: logics.TASK_STATUS_被收录})
		if includeNum == logics.INCLUDE_有收录却无收录数 {
			days := services.DaysApartToday(finishedTask.EnginedDomain.CheckedAt)
			if days > logics.INCLUDE_查询有收录却无收录数站点的相隔天数 {
				databases.Db.Model(&entities.KeywordInclude{}).Where("engined_domain_id = ?", finishedTask.EnginedDomain.ID).Update(map[string]interface{}{
					"need_check": true,
				})
			}
		} else if includeNum != oldIncludeNum {
			databases.Db.Model(&entities.KeywordInclude{}).Where("engined_domain_id = ?", finishedTask.EnginedDomain.ID).Update(map[string]interface{}{
				"need_check": true,
			})
		}
	}
}
