package keyword_service

import (
	"app/databases"
	"app/databases/entities"
)

func IncludedKeywords(domain string, engine string) ([]string, error) {
	var enginedDomain entities.EnginedDomain
	databases.Db.Model(entities.EnginedDomain{}).
		Where("domain = ?", domain).
		Where("engine = ?", engine).
		First(&enginedDomain)
	result := make([]string, 0)
	err := databases.Db.Model(entities.KeywordInclude{}).
		Where("engined_domain_id = ?", enginedDomain.ID).
		Where("is_included = ?", IS_INCLUDE_YES).
		Pluck("keyword", &result).Error
	return result, err
}
