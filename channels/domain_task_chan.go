package channels

import (
	"app/databases/entities"
	"app/structs/models/logics"
)

var DomainTaskSendingChan chan *entities.DomainIncludeTask

func init() {
	DomainTaskSendingChan = make(chan *entities.DomainIncludeTask, logics.TASK_发送下载缓冲区大小)
}
