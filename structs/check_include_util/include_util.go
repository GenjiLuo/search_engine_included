package check_include_util

type RequestBuilderRequest struct {
	SearchWord  string `json:"search_word" bson:"search_word"`
	Page        int    `json:"page" bson:"page"`
	Capture     bool   `json:"capture" bson:"capture"`
	SearchCycle int    `json:"search_cycle"`
	Priority    string `json:"priority"`
}

type ParseIncludeRequest struct {
	Body string `json:"body" bson:"body"`
}

type ParseIncludeResponse struct {
	IncludeNum int `json:"include_num"`
}

type KeywordParseIncludeResponse struct {
	IsIncluded bool `json:"is_included"`
}
