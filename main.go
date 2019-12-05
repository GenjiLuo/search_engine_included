package main

import (
	"app/databases"
	"app/jobs"
	"time"
)

func foreverGo(run func(), routineLimits int) {
	for i := 0; i < routineLimits; i++ {
		go func() {
			for {
				run()
			}
		}()
	}
}

func main() {
	databases.AutoMigrate()

	foreverGo(jobs.AddDomainsToTask, 1)
	foreverGo(jobs.SendDomainTasksToChan, 1)
	foreverGo(jobs.SendDomainDcRequestsFromChan, 5)
	foreverGo(jobs.FetchDomainQueryParsedResult, 1)
	foreverGo(jobs.AddKeywordsToTask, 1)
	foreverGo(jobs.SendKeywordTasksToChan, 1)
	foreverGo(jobs.SendKeywordDcRequestsFromChan, 10)
	foreverGo(jobs.FetchKeywordQueryParsedResult, 1)
	foreverGo(jobs.SyncRankTasks, 1)

	for {
		time.Sleep(time.Minute)
	}
}
