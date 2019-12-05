package entities

import "time"

type EnginedDomain struct {
	ID              int       `gorm:"primary_key" json:"id"`
	Domain          string    `gorm:"type:varchar(255);unique_index:domain_engine" json:"domain"`
	Engine          string    `gorm:"type:varchar(255);unique_index:domain_engine" json:"engine"`
	IncludeNum      int       `gorm:"type:int;" json:"include_num"`
	CheckedAt       time.Time `json:"checked_at"`
	KeywordIncludes []KeywordInclude
}

type KeywordInclude struct {
	ID              int       `gorm:"primary_key" json:"id"`
	EnginedDomainId int       `gorm:"type:int;index:search;unique_index:search_keyword" json:"engined_domain_id"`
	Keyword         string    `gorm:"type:varchar(255);index:search;unique_index:search_keyword" json:"keyword"`
	NeedCheck       bool      `gorm:"type:tinyint;default:0;index:need_check;" json:"need_check"`
	NoIncludedDays  int       `gorm:"type:int;default:0" json:"no_included_days"`
	IsIncluded      bool      `gorm:"type:tinyint;default:1;index:is_included;" json:"is_included"`
	LastIncludedAt  time.Time `json:"last_included_at"`
	EnginedDomain   EnginedDomain
}

type DomainIncludeTask struct {
	ID              int    `gorm:"primary_key" json:"id"`
	EnginedDomainId int    `gorm:"type:int;index:search" json:"engined_domain_id"`
	UniqueKey       string `gorm:"type:varchar(255);index:unique_key;" json:"unique_key"`
	Status          int    `gorm:"type:int;index:task_status" json:"status"`
	CreatedAt       time.Time
	EnginedDomain   EnginedDomain
}

type KeywordIncludeTask struct {
	ID               int    `gorm:"primary_key" json:"id"`
	KeywordIncludeId int    `gorm:"type:int;" json:"keyword_include_id"`
	UniqueKey        string `gorm:"type:varchar(255);index:unique_key;" json:"unique_key"`
	Status           int    `gorm:"type:int;index:task_status" json:"status"`
	CreatedAt        time.Time
	KeywordInclude   KeywordInclude
}
