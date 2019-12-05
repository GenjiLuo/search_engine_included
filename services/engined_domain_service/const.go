package engined_domain_service

const ENGINE_BAIDU_PC = "baidu_pc"
const ENGINE_BAIDU_MOBILE = "baidu_mobile"
const ENGINE_SOGOU_PC = "sogou_pc"
const ENGINE_SM_MOBILE = "sm_mobile"
const ENGINE_360_PC = "360_pc"

func Engines() []string {
	return []string{
		ENGINE_BAIDU_PC,
		ENGINE_BAIDU_MOBILE,
		ENGINE_SOGOU_PC,
		ENGINE_SM_MOBILE,
		ENGINE_360_PC,
	}
}
