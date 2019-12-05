package keyword_service

import (
	"app/databases"
	"app/services/engined_domain_service"
	"fmt"
	"strings"
	"time"
)

func KeywordsAdd(domain string, keywords []string) error {
	engines := engined_domain_service.Engines()
	for _, engine := range engines {
		enginedDomain, err := engined_domain_service.EnsureEnginedDomain(domain, engine)
		if err != nil {
			return err
		}

		valueStrings := make([]string, 0)
		valueArgs := make([]interface{}, 0)
		for _, keyword := range keywords {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, enginedDomain.ID)
			valueArgs = append(valueArgs, keyword)
			valueArgs = append(valueArgs, time.Now())
		}

		if len(valueStrings) > 0 {
			sql := fmt.Sprintf("INSERT IGNORE INTO `keyword_include` (`engined_domain_id`, `keyword`, `last_included_at`) VALUES %s", strings.Join(valueStrings, ","))
			err := databases.Db.Exec(sql, valueArgs...).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
