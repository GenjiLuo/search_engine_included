package databases

import (
	"app/settings"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var Db *gorm.DB

func init() {
	var err error

	connArgs := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.DbUsername,
		settings.DbPassword,
		settings.DbHost,
		settings.DbPort,
		settings.DbDatabase,
	)

	Db, err = gorm.Open(settings.DbConnection, connArgs)

	if err != nil {
		log.Fatalln(err)
	}

	Db.SingularTable(true)
	Db = Db.Set("gorm:association_autoupdate", false)

	Db.DB().SetMaxOpenConns(100)
	Db.DB().SetMaxIdleConns(20)
	Db.DB().SetConnMaxLifetime(55 * time.Second)

	if settings.Debug {
		Db = Db.Debug()
	}
}
