package keyword_service

import (
	"app/databases"
	"app/databases/entities"
	"app/services/engined_domain_service"
)

func KeywordsRemove(domain string, keywords []string) error {
	engines := engined_domain_service.Engines()
	enginedDomainIds := make([]int, 0)
	databases.Db.Model(entities.EnginedDomain{}).
		Where(entities.EnginedDomain{Domain: domain}).
		Where("engine in (?)", engines).
		Pluck("id", &enginedDomainIds)

	err := databases.Db.Model(entities.KeywordInclude{}).
		Where("engined_domain_id in (?)", enginedDomainIds).
		Where("keyword in (?)", keywords).
		Delete(entities.KeywordInclude{}).Error
	if err != nil {
		return err
	}

	return nil
}
