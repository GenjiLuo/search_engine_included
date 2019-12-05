package databases

import "app/databases/entities"

func AutoMigrate() {
	Db.AutoMigrate(
		&entities.EnginedDomain{},
		&entities.KeywordInclude{},
		&entities.DomainIncludeTask{},
		&entities.KeywordIncludeTask{},
	)
}
