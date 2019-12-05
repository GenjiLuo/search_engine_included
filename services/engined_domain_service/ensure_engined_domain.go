package engined_domain_service

import (
	"app/databases"
	"app/databases/entities"
	"time"
)

func EnsureEnginedDomain(domain, engine string) (*entities.EnginedDomain, error) {
	enginedDomain := entities.EnginedDomain{
		Domain: domain,
		Engine: engine,
	}
	err := databases.Db.Where(enginedDomain).
		Attrs(entities.EnginedDomain{
			CheckedAt: time.Now(),
		}).
		FirstOrCreate(&enginedDomain).Error
	if err != nil {
		return nil, err
	}

	return &enginedDomain, nil
}
