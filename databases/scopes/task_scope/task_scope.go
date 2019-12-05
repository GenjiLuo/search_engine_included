package task_scope

import (
	"app/structs/models/logics"
	"github.com/jinzhu/gorm"
)

func UnQueried(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", logics.TASK_STATUS_未查询)
}

func Querying(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", logics.TASK_STATUS_查询中)
}

func BeforeQueried(db *gorm.DB) *gorm.DB {
	return db.Where("status IN (?)", []int{logics.TASK_STATUS_未查询, logics.TASK_STATUS_查询中})
}

func UniqueKeysIn(uniqueKeys []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("unique_key IN (?)", uniqueKeys)
	}
}
