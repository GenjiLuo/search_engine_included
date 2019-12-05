package task_service_test

import (
	"app/databases/entities"
	"app/services/task_service"
	"testing"
)

func TestKeywordUniqueKeyMappedTasks(t *testing.T) {
	tasks := []*entities.KeywordIncludeTask{
		{ID: 1, UniqueKey: "abc"},
		{ID: 2, UniqueKey: "abc"},
		{ID: 3, UniqueKey: "bcd"},
	}
	mappedTasks := task_service.KeywordUniqueKeyMappedTasks(tasks)
	abcTasks := mappedTasks["abc"]
	if abcTasks[0].ID != 1 || abcTasks[1].ID != 2 {
		t.Errorf("abc ID error, expect %d, %d, got %d, %d", tasks[0].ID, tasks[1].ID, abcTasks[0].ID, abcTasks[1].ID)
	}
	bcdTasks := mappedTasks["bcd"]
	if bcdTasks[0].ID != 3 {
		t.Errorf("bcd ID error, expect %d, got %d", tasks[2].ID, bcdTasks[0].ID)
	}
}
