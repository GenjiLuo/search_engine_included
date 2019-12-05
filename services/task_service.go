package services

import (
	"app/databases"
	"app/databases/db_keyword_include_service"
	"app/databases/entities"
	"app/structs/models/logics"
	"time"
)

// AddSomeUnCheckDomainToTasks: 按序添加一些主域到任务
func AddSomeUnCheckDomainToTasks() int {
	var domainIds []int
	databases.Db.Table("engined_domain").
		Joins("left join domain_include_task on engined_domain.id=domain_include_task.engined_domain_id ").
		Where("domain_include_task.id is NULL").
		Order("checked_at").
		Limit(logics.TASK_单次放入任务数量限制).
		Pluck("engined_domain.id", &domainIds)

	if len(domainIds) != 0 {
		databases.Db.Exec("INSERT INTO domain_include_task (`engined_domain_id`, `status`, `created_at`) (SELECT id, 1, NOW() FROM engined_domain WHERE id in (?))", domainIds)
	}
	return len(domainIds)
}

// AddSomeUnCheckKeywordToTasks: 添加一些关键词到任务
func AddSomeUnCheckKeywordToTasks() int {
	var keywordIds []int
	databases.Db.Model(&entities.KeywordInclude{}).
		Where("need_check = ?", 1).
		Limit(logics.TASK_单次放入任务数量限制).
		Pluck("id", &keywordIds)

	if len(keywordIds) != 0 {
		db_keyword_include_service.NotNeedCheck(keywordIds)
		databases.Db.Exec("INSERT INTO keyword_include_task (`keyword_include_id`, `status`, `created_at`) (SELECT id, 1, NOW() FROM keyword_include WHERE id in (?))", keywordIds)
	}
	return len(keywordIds)
}

// DaysApartToday: 某个时间点距离今天相隔的天数
func DaysApartToday(formerDay time.Time) int {
	t, _ := time.ParseInLocation("2006-01-02", formerDay.Format("2006-01-02"), time.Local)
	duration := time.Now().Sub(t)
	return int(duration.Hours() / 24)
}
