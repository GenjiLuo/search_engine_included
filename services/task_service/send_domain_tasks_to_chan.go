package task_service

import "app/databases/entities"

func SendDomainTasksToChan(domainTasks []*entities.DomainIncludeTask, domainTaskChan chan *entities.DomainIncludeTask, preFunc func(domainTask *entities.DomainIncludeTask)) {
	for _, domainTask := range domainTasks {
		preFunc(domainTask)
		domainTaskChan <- domainTask
	}
}
