package download_center

type DcRequest struct {
	UniqueKey string `json:"unique_key"`
	Request   struct {
		Url       string `json:"url"`
		UserAgent string `json:"user_agent"`
		Cookie    string `json:"cookie"`
		Body      string `json:"body"`
	} `json:"request"`
	Config struct {
		District       string   `json:"district"`
		ResponseTypes  []string `json:"response_types" bson:"response_types"`
		FollowRedirect bool     `json:"follow_redirect"`
		Priority       string   `json:"priority"`
	} `json:"config"`
	Status int `json:"status"`
}

type DcResponse struct {
	Header  string `json:"header"`
	Body    string `json:"body"`
}
