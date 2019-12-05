package task_service_test

import (
	"app/databases/entities"
	"app/services/task_service"
	"app/structs/models/logics"
	"testing"
)

func TestSendDomainTasksToChan(t *testing.T) {
	var tasks []*entities.DomainIncludeTask
	task1 := &entities.DomainIncludeTask{EnginedDomainId: 1}
	task2 := &entities.DomainIncludeTask{EnginedDomainId: 2}
	tasks = append(tasks, task1, task2)
	channel := make(chan *entities.DomainIncludeTask, logics.TASK_发送下载缓冲区大小)

	task_service.SendDomainTasksToChan(tasks, channel, func(task *entities.DomainIncludeTask) {
		task.Status = logics.TASK_STATUS_查询中
	})

	verifyTask1 := <-channel
	verifyTask2 := <-channel

	if verifyTask1.EnginedDomainId != verifyTask1.EnginedDomainId {
		t.Error("keywordId not match")
	}
	if verifyTask2.EnginedDomainId != verifyTask2.EnginedDomainId {
		t.Error("keywordId not match")
	}
	if verifyTask1.Status != logics.TASK_STATUS_查询中 {
		t.Error("status not changed")
	}
	if verifyTask2.Status != logics.TASK_STATUS_查询中 {
		t.Error("status not changed")
	}
}
