package task_service_test

import (
	"app/databases/entities"
	"app/services/task_service"
	"app/structs/models/logics"
	"testing"
)

func TestSendKeywordTasksToChan(t *testing.T) {
	var tasks []*entities.KeywordIncludeTask
	task1 := &entities.KeywordIncludeTask{KeywordIncludeId: 1}
	task2 := &entities.KeywordIncludeTask{KeywordIncludeId: 2}
	tasks = append(tasks, task1, task2)
	channel := make(chan *entities.KeywordIncludeTask, logics.TASK_发送下载缓冲区大小)

	task_service.SendKeywordTasksToChan(tasks, channel, func(task *entities.KeywordIncludeTask) {
		task.Status = logics.TASK_STATUS_查询中
	})

	verifyTask1 := <-channel
	verifyTask2 := <-channel

	if verifyTask1.KeywordIncludeId != verifyTask1.KeywordIncludeId {
		t.Error("keywordId not match")
	}
	if verifyTask2.KeywordIncludeId != verifyTask2.KeywordIncludeId {
		t.Error("keywordId not match")
	}
	if verifyTask1.Status != logics.TASK_STATUS_查询中 {
		t.Error("status not changed")
	}
	if verifyTask2.Status != logics.TASK_STATUS_查询中 {
		t.Error("status not changed")
	}
}
