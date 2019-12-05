package task_service

import "app/databases/entities"

func KeywordUniqueKeyMappedTasks(finishedTasks []*entities.KeywordIncludeTask) map[string][]*entities.KeywordIncludeTask {
	uniqueKeyTasksMap := make(map[string][]*entities.KeywordIncludeTask)
	for i, _ := range finishedTasks {
		uniqueKeyTasksMap[finishedTasks[i].UniqueKey] = append(uniqueKeyTasksMap[finishedTasks[i].UniqueKey], finishedTasks[i])
	}

	return uniqueKeyTasksMap
}
