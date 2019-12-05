package channels

import (
	"app/databases/entities"
	"app/structs/models/logics"
)

var KeywordTaskSendingChan chan *entities.KeywordIncludeTask

func init() {
	KeywordTaskSendingChan = make(chan *entities.KeywordIncludeTask, logics.TASK_发送下载缓冲区大小)
}
