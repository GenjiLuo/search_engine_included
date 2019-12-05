package check_include_util

import "app/services/engined_domain_service"

var MapEngine map[string]string

func init() {
	MapEngine = map[string]string{
		engined_domain_service.ENGINE_BAIDU_PC:     "baidu-pc",
		engined_domain_service.ENGINE_BAIDU_MOBILE: "baidu-pc",
		engined_domain_service.ENGINE_SOGOU_PC:     "sogou-pc",
		engined_domain_service.ENGINE_360_PC:       "360-pc",
		engined_domain_service.ENGINE_SM_MOBILE:    "sm-mobile",
	}
}
