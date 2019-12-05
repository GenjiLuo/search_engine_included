package task_service

import "app/databases/entities"

func UniqueKeyMappedTasks(finishedTasks []*entities.DomainIncludeTask) map[string][]*entities.DomainIncludeTask {
	uniqueKeyTasksMap := make(map[string][]*entities.DomainIncludeTask)
	for i, _ := range finishedTasks {
		uniqueKeyTasksMap[finishedTasks[i].UniqueKey] = append(uniqueKeyTasksMap[finishedTasks[i].UniqueKey], finishedTasks[i])
	}

	return uniqueKeyTasksMap
}
