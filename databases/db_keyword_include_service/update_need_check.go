package db_keyword_include_service

import (
	"app/databases"
	"app/databases/entities"
	"github.com/jinzhu/gorm"
)

//NotNeedCheck 设置need_check为不查
func NotNeedCheck(ids []int) *gorm.DB {
	if len(ids) == 0 {
		return nil
	}

	return databases.Db.Model(entities.KeywordInclude{}).
		Where("id in (?)", ids).
		Updates(map[string]interface{}{"need_check": false})
}

//NeedCheck 设置need_check为查
func NeedCheck(ids []int) *gorm.DB {
	if len(ids) == 0 {
		return nil
	}

	return databases.Db.Model(&entities.KeywordInclude{}).
		Where("id in (?)", ids).
		Updates(map[string]interface{}{"need_check": true})
}
