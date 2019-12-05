package task_service

import "app/databases/entities"

func SendKeywordTasksToChan(keywordTasks []*entities.KeywordIncludeTask, keywordTaskChan chan *entities.KeywordIncludeTask, preFunc func(keywordTask *entities.KeywordIncludeTask)) {
	for _, keywordTask := range keywordTasks {
		preFunc(keywordTask)
		keywordTaskChan <- keywordTask
	}
}
