package pkg

const (
	RouteTypeAdmin    = "admin"
	RouteTypeIntranet = "intranet"
	RouteTypeOuternet = "outernet"
	RouteTypePage     = "page"
)
const (
	EnvDev  = "dev"  // 开发环境
	EnvTest = "test" // 测试环境
	EnvPre  = "pre"  // 预发布环境
	EnvProd = "prod" // 线上环境
)
const (
	GatewayErrorCodeSign              = iota + 10000 // 签名验证失败
	GatewayErrorCodeNotLogin                         // 未登录
	GatewayErrorCodeNotHavePermit                    // 无权限
	GatewayErrorCodePermitConfigError                // 网关配置错误
)

var (
	MapGatewayError = map[int]string{
		GatewayErrorCodeSign:              "签名验证失败",
		GatewayErrorCodeNotLogin:          "未登录",
		GatewayErrorCodeNotHavePermit:     "没有权限",
		GatewayErrorCodePermitConfigError: "网关配置错误(验证权限必须要验证用户登录%s)",
	}
)

type ResponseCallBack func()
